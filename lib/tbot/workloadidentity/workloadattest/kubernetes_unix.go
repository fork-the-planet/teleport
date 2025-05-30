//go:build unix

/*
 * Teleport
 * Copyright (C) 2024  Gravitational, Inc.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package workloadattest

import (
	"cmp"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"
	v1 "k8s.io/api/core/v1"

	workloadidentityv1pb "github.com/gravitational/teleport/api/gen/proto/go/teleport/workloadidentity/v1"
	"github.com/gravitational/teleport/api/utils/retryutils"
	"github.com/gravitational/teleport/lib/tbot/workloadidentity/workloadattest/container"
)

// imageDigestRegex finds the `sha256:<hash>` image digest in an OCI image URI.
//
// We use a regex rather than parsing the URI because the format may differ
// between the schemes, and it's difficult to find comprehensive documentation.
var imageDigestRegex = regexp.MustCompile(`sha256:([[:xdigit:]]{64})`)

// KubernetesAttestor attests a workload to a Kubernetes pod.
//
// It requires:
//
// - `hostPID: true` so we can view the /proc of other pods.
// - `TELEPORT_MY_NODE_NAME` to be set to the node name of the current node.
// - A service account that allows it to query the Kubelet API.
//
// It roughly takes the following steps:
//  1. From the PID, determine the container ID and pod ID from the
//     /proc/<pid>/mountinfo file.
//  2. Makes a request to the Kubelet API to list all pods on the node.
//  3. Find the pod and container with the matching ID.
//  4. Convert the pod information to a KubernetesAttestation.
type KubernetesAttestor struct {
	kubeletClient *kubeletClient
	log           *slog.Logger
	// rootPath specifies the location of `/`. This allows overriding for tests.
	rootPath string
	clock    clockwork.Clock
}

// NewKubernetesAttestor creates a new KubernetesAttestor.
func NewKubernetesAttestor(cfg KubernetesAttestorConfig, log *slog.Logger) *KubernetesAttestor {
	kubeletClient := newKubeletClient(cfg.Kubelet)
	return &KubernetesAttestor{
		kubeletClient: kubeletClient,
		log:           log,
		clock:         clockwork.NewRealClock(),
	}
}

// Attest resolves the Kubernetes pod information from the
// PID of the workload.
func (a *KubernetesAttestor) Attest(ctx context.Context, pid int) (*workloadidentityv1pb.WorkloadAttrsKubernetes, error) {
	a.log.DebugContext(ctx, "Starting Kubernetes workload attestation", "pid", pid)

	container, err := container.LookupPID(a.rootPath, pid, container.KubernetesParser)
	if err != nil {
		return nil, trace.Wrap(err, "determining pod and container ID")
	}

	a.log.DebugContext(ctx,
		"Found pod and container ID",
		"pod_id", container.PodID,
		"container_id", container.ID,
	)

	pod, containerStatus, err := a.getPodAndContainerStatus(ctx, container.PodID, container.ID)
	if err != nil {
		return nil, trace.Wrap(err, "finding pod by ID")
	}
	a.log.DebugContext(ctx, "Found pod", "pod_name", pod.Name)

	var ctr *workloadidentityv1pb.WorkloadAttrsKubernetesContainer
	if containerStatus != nil {
		ctr = &workloadidentityv1pb.WorkloadAttrsKubernetesContainer{
			Name:        containerStatus.Name,
			Image:       containerStatus.Image,
			ImageDigest: imageDigestRegex.FindString(containerStatus.ImageID),
		}
	}

	att := &workloadidentityv1pb.WorkloadAttrsKubernetes{
		Attested:       true,
		Namespace:      pod.Namespace,
		ServiceAccount: pod.Spec.ServiceAccountName,
		PodName:        pod.Name,
		PodUid:         string(pod.UID),
		Labels:         pod.Labels,
		Container:      ctr,
	}
	a.log.DebugContext(ctx, "Finished Kubernetes workload attestation", "attestation", att)
	return att, nil
}

func (a *KubernetesAttestor) getPodAndContainerStatus(ctx context.Context, podID, containerID string) (*v1.Pod, *v1.ContainerStatus, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	log := a.log.With("pod_id", podID, "container_id", containerID)

	retry, err := retryutils.NewRetryV2(retryutils.RetryV2Config{
		Driver: retryutils.NewExponentialDriver(100 * time.Millisecond),
		Max:    2 * time.Second,
		Clock:  a.clock,
	})
	if err != nil {
		return nil, nil, trace.Wrap(err, "creating retrier")
	}

	var (
		pod             *v1.Pod
		containerStatus *v1.ContainerStatus
	)
LOOP:
	for {
		pod, containerStatus, err = a.tryGetPodAndContainerStatus(ctx, podID, containerID)
		switch {
		case err != nil:
			return nil, nil, err
		case containerStatus == nil:
			// It's possible for a workload container to start and request a SVID
			// before the kubelet has updated its state, in which case we might
			// get back no container status at all, or in the case of a restart,
			// the previous run's status.
			log.DebugContext(ctx, "Kubelet did not return expected container status; its state might be stale")
		default:
			break LOOP
		}

		retry.Inc()
		select {
		case <-ctx.Done():
			break LOOP
		case <-retry.After():
		}
	}

	if pod != nil {
		return pod, containerStatus, nil
	}
	return nil, nil, err
}

func (a *KubernetesAttestor) tryGetPodAndContainerStatus(ctx context.Context, podID, containerID string) (*v1.Pod, *v1.ContainerStatus, error) {
	pods, err := a.kubeletClient.ListAllPods(ctx)
	if err != nil {
		return nil, nil, trace.Wrap(err, "listing all pods")
	}

	var pod *v1.Pod
	for _, p := range pods.Items {
		if string(p.UID) == podID {
			pod = &p
			break
		}
	}
	if pod == nil {
		return nil, nil, trace.NotFound("pod %q not found", podID)
	}

	var containerStatus *v1.ContainerStatus
	for _, status := range pod.Status.ContainerStatuses {
		// Kubelet returns the container ID prefixed by `<type>://`.
		if _, id, _ := strings.Cut(status.ContainerID, "://"); id == containerID {
			containerStatus = &status
			break
		}
	}
	return pod, containerStatus, nil
}

// kubeletClient is a HTTP client for the Kubelet API
type kubeletClient struct {
	cfg    KubeletClientConfig
	getEnv func(string) string
}

func newKubeletClient(cfg KubeletClientConfig) *kubeletClient {
	return &kubeletClient{
		cfg:    cfg,
		getEnv: os.Getenv,
	}
}

type roundTripperFn func(req *http.Request) (*http.Response, error)

func (f roundTripperFn) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func (c *kubeletClient) httpClient() (url.URL, *http.Client, error) {
	host := c.getEnv(nodeNameEnv)

	if c.cfg.ReadOnlyPort != 0 {
		return url.URL{
			Scheme: "http",
			Host:   net.JoinHostPort(host, strconv.Itoa(c.cfg.ReadOnlyPort)),
		}, &http.Client{}, nil
	}

	port := cmp.Or(c.cfg.SecurePort, defaultSecurePort)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}

	switch {
	case c.cfg.SkipVerify:
		transport.TLSClientConfig.InsecureSkipVerify = true
	default:
		caPath := cmp.Or(c.cfg.CAPath, defaultCAPath)
		certPool := x509.NewCertPool()
		caPEM, err := os.ReadFile(caPath)
		if err != nil {
			return url.URL{}, nil, trace.Wrap(err, "reading CA file %q", caPath)
		}
		if !certPool.AppendCertsFromPEM(caPEM) {
			return url.URL{}, nil, trace.BadParameter("failed to append CA cert from %q", caPath)
		}
		transport.TLSClientConfig.RootCAs = certPool
	}

	client := &http.Client{
		Transport: transport,
		// 10 seconds is fairly generous given that we're expecting to talk to
		// kubelet on the same physical machine.
		Timeout: 10 * time.Second,
	}

	switch {
	case c.cfg.Anonymous:
	// Nothing to do
	case c.cfg.TokenPath != "":
		fallthrough
	default:
		tokenPath := cmp.Or(c.cfg.TokenPath, defaultServiceAccountTokenPath)
		token, err := os.ReadFile(tokenPath)
		if err != nil {
			return url.URL{}, nil, trace.Wrap(err, "reading token file %q", tokenPath)
		}
		client.Transport = roundTripperFn(func(req *http.Request) (*http.Response, error) {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			return transport.RoundTrip(req)
		})
	}

	return url.URL{
		Scheme: "https",
		Host:   net.JoinHostPort(host, strconv.Itoa(port)),
	}, client, nil
}

func (c *kubeletClient) ListAllPods(ctx context.Context) (*v1.PodList, error) {
	reqUrl, client, err := c.httpClient()
	if err != nil {
		return nil, trace.Wrap(err, "creating HTTP client")
	}
	reqUrl.Path = "/pods"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl.String(), nil)
	if err != nil {
		return nil, trace.Wrap(err, "creating request")
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, trace.Wrap(err, "performing request")
	}
	defer res.Body.Close()

	out := &v1.PodList{}
	if err := json.NewDecoder(res.Body).Decode(out); err != nil {
		return nil, trace.Wrap(err, "decoding response")
	}
	return out, nil
}
