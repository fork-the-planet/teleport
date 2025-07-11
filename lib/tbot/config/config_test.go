/*
 * Teleport
 * Copyright (C) 2023  Gravitational, Inc.
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

package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/gravitational/teleport/lib/defaults"
	"github.com/gravitational/teleport/lib/tbot/bot"
	"github.com/gravitational/teleport/lib/tbot/botfs"
	"github.com/gravitational/teleport/lib/utils/testutils/golden"
)

func TestConfigFile(t *testing.T) {
	configData := fmt.Sprintf(exampleConfigFile, "foo")
	cfg, err := ReadConfig(strings.NewReader(configData), false)
	require.NoError(t, err)
	require.NoError(t, cfg.CheckAndSetDefaults())

	require.Equal(t, "auth.example.com", cfg.AuthServer)
	require.Equal(t, time.Minute*5, cfg.CredentialLifetime.RenewalInterval)

	require.NotNil(t, cfg.Onboarding)

	token, err := cfg.Onboarding.Token()
	require.NoError(t, err)
	require.Equal(t, "foo", token)
	require.ElementsMatch(t, []string{"sha256:abc123"}, cfg.Onboarding.CAPins)

	_, ok := cfg.Storage.Destination.(*DestinationMemory)
	require.True(t, ok)

	require.Len(t, cfg.Services, 1)
	output := cfg.Services[0]
	identOutput, ok := output.(*IdentityOutput)
	require.True(t, ok)

	destImpl := identOutput.GetDestination()
	destImplReal, ok := destImpl.(*DestinationDirectory)
	require.True(t, ok)
	require.Equal(t, "/tmp/foo", destImplReal.Path)

	require.True(t, cfg.Debug)
	require.Equal(t, "127.0.0.1:1337", cfg.DiagAddr)
}

func TestLoadTokenFromFile(t *testing.T) {
	tokenDir := t.TempDir()
	tokenFile := filepath.Join(tokenDir, "token")
	require.NoError(t, os.WriteFile(tokenFile, []byte("xxxyyy"), 0660))

	configData := fmt.Sprintf(exampleConfigFile, tokenFile)
	cfg, err := ReadConfig(strings.NewReader(configData), false)
	require.NoError(t, err)

	token, err := cfg.Onboarding.Token()
	require.NoError(t, err)
	require.Equal(t, "xxxyyy", token)
}

const exampleConfigFile = `
version: v2
auth_server: auth.example.com
renewal_interval: 5m
debug: true
diag_addr: 127.0.0.1:1337
onboarding:
  token: %s
  ca_pins:
    - sha256:abc123
storage:
  type: memory
outputs:
  - type: identity
    destination:
      type: directory
      path: /tmp/foo
`

func TestDestinationFromURI(t *testing.T) {
	tests := []struct {
		in      string
		want    bot.Destination
		wantErr bool
	}{
		{
			in: "/absolute/dir",
			want: &DestinationDirectory{
				Path: "/absolute/dir",
			},
		},
		{
			in: "relative/dir",
			want: &DestinationDirectory{
				Path: "relative/dir",
			},
		},
		{
			in: "./relative/dir",
			want: &DestinationDirectory{
				Path: "./relative/dir",
			},
		},
		{
			in: "file:///absolute/dir",
			want: &DestinationDirectory{
				Path: "/absolute/dir",
			},
		},
		{
			in: "file:/absolute/dir",
			want: &DestinationDirectory{
				Path: "/absolute/dir",
			},
		},
		{
			in:      "file://host/absolute/dir",
			wantErr: true,
		},
		{
			in:   "memory://",
			want: &DestinationMemory{},
		},
		{
			in:      "memory://foo/bar",
			wantErr: true,
		},
		{
			in:      "foobar://",
			wantErr: true,
		},
		{
			in: "kubernetes-secret:///my-secret",
			want: &DestinationKubernetesSecret{
				Name: "my-secret",
			},
		},
		{
			in: "kubernetes-secret://my-secret",
			want: &DestinationKubernetesSecret{
				Name: "my-secret",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := DestinationFromURI(tt.in)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

// TestBotConfig_YAML ensures that as a whole YAML marshaling and unmarshaling
// of the config works as expected. Avoid testing exhaustive cases here and
// prefer the Output YAML tests for testing the intricacies of marshaling and
// unmarshaling specific objects.
func TestBotConfig_YAML(t *testing.T) {
	tests := []testYAMLCase[BotConfig]{
		{
			name: "standard config",
			in: BotConfig{
				Version: V2,
				Storage: &StorageConfig{
					Destination: &DestinationDirectory{
						Path:     "/bot/storage",
						ACLs:     botfs.ACLTry,
						Symlinks: botfs.SymlinksSecure,
					},
				},
				Onboarding: OnboardingConfig{
					JoinMethod: "gitlab",
					TokenValue: "my-gitlab-token",
					Gitlab: GitlabOnboardingConfig{
						TokenEnvVarName: "MY_CUSTOM_ENV_VAR",
					},
				},
				FIPS:       true,
				Debug:      true,
				Oneshot:    true,
				AuthServer: "example.teleport.sh:443",
				DiagAddr:   "127.0.0.1:1337",
				CredentialLifetime: CredentialLifetime{
					TTL:             time.Minute,
					RenewalInterval: time.Second * 30,
				},
				Outputs: ServiceConfigs{
					&IdentityOutput{
						Destination: &DestinationDirectory{
							Path: "/bot/output",
						},
						Roles:   []string{"editor"},
						Cluster: "example.teleport.sh",
					},
					&IdentityOutput{
						Destination: &DestinationMemory{},
					},
					&IdentityOutput{
						Destination: &DestinationKubernetesSecret{
							Name: "my-secret",
						},
						CredentialLifetime: CredentialLifetime{
							TTL:             30 * time.Second,
							RenewalInterval: 15 * time.Second,
						},
					},
				},
				Services: []ServiceConfig{
					&SPIFFEWorkloadAPIService{
						Listen: "unix:///var/run/spiffe.sock",
						SVIDs: []SVIDRequestWithRules{
							{
								SVIDRequest: SVIDRequest{
									Path: "/bar",
									Hint: "my hint",
									SANS: SVIDRequestSANs{
										DNS: []string{"foo.bar"},
										IP:  []string{"10.0.0.1"},
									},
								},
								Rules: []SVIDRequestRule{
									{
										Unix: SVIDRequestRuleUnix{
											PID: ptr(100),
											UID: ptr(1000),
											GID: ptr(1234),
										},
									},
									{
										Unix: SVIDRequestRuleUnix{
											PID: ptr(100),
										},
									},
								},
							},
						},
						CredentialLifetime: CredentialLifetime{
							TTL:             30 * time.Second,
							RenewalInterval: 15 * time.Second,
						},
					},
					&ExampleService{
						Message: "llama",
					},
					&SSHMultiplexerService{
						Destination: &DestinationDirectory{
							Path: "/bot/output",
						},
						CredentialLifetime: CredentialLifetime{
							TTL:             30 * time.Second,
							RenewalInterval: 15 * time.Second,
						},
					},
					&ApplicationTunnelService{
						Listen:  "tcp://127.0.0.1:123",
						Roles:   []string{"access"},
						AppName: "my-app",
						CredentialLifetime: CredentialLifetime{
							TTL:             30 * time.Second,
							RenewalInterval: 15 * time.Second,
						},
					},
					&WorkloadIdentityX509Service{
						Destination: &DestinationDirectory{
							Path: "/an/output/path",
						},
						Selector: WorkloadIdentitySelector{
							Name: "my-workload-identity",
						},
						CredentialLifetime: CredentialLifetime{
							TTL:             30 * time.Second,
							RenewalInterval: 15 * time.Second,
						},
					},
					&WorkloadIdentityAPIService{
						Listen: "tcp://127.0.0.1:123",
						Selector: WorkloadIdentitySelector{
							Name: "my-workload-identity",
						},
						CredentialLifetime: CredentialLifetime{
							TTL:             30 * time.Second,
							RenewalInterval: 15 * time.Second,
						},
					},
					&WorkloadIdentityJWTService{
						Destination: &DestinationDirectory{
							Path: "/an/output/path",
						},
						Selector: WorkloadIdentitySelector{
							Name: "my-workload-identity",
						},
						Audiences: []string{"audience1", "audience2"},
					},
				},
			},
		},
		{
			name: "minimal config",
			in: BotConfig{
				Version:    V2,
				AuthServer: "example.teleport.sh:443",
				CredentialLifetime: CredentialLifetime{
					TTL:             time.Minute,
					RenewalInterval: time.Second * 30,
				},
				Outputs: ServiceConfigs{
					&IdentityOutput{
						Destination: &DestinationMemory{},
					},
				},
			},
		},
		{
			name: "minimal config using proxy addr",
			in: BotConfig{
				Version:     V2,
				ProxyServer: "example.teleport.sh:443",
				CredentialLifetime: CredentialLifetime{
					TTL:             time.Minute,
					RenewalInterval: time.Second * 30,
				},
				Outputs: ServiceConfigs{
					&IdentityOutput{
						Destination: &DestinationMemory{},
					},
				},
			},
		},
	}

	testYAML(t, tests)
}

type testYAMLCase[T any] struct {
	name string
	in   T
}

func testYAML[T any](t *testing.T, tests []testYAMLCase[T]) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBuffer(nil)
			encoder := yaml.NewEncoder(b)
			encoder.SetIndent(2)
			require.NoError(t, encoder.Encode(&tt.in))

			if golden.ShouldSet() {
				golden.Set(t, b.Bytes())
			}
			require.Equal(
				t,
				string(golden.Get(t)),
				b.String(),
				"results of marshal did not match golden file, rerun tests with GOLDEN_UPDATE=1",
			)

			// Now test unmarshalling to see if we get the same object back
			decoder := yaml.NewDecoder(b)
			var unmarshalled T
			require.NoError(t, decoder.Decode(&unmarshalled))
			require.Equal(t, tt.in, unmarshalled, "unmarshalling did not result in same object as input")
		})
	}
}

func TestBotConfig_InsecureWithCAPins(t *testing.T) {
	cfg := &BotConfig{
		Insecure: true,
		Onboarding: OnboardingConfig{
			CAPins: []string{"123"},
		},
	}

	require.ErrorContains(t, cfg.CheckAndSetDefaults(), "ca-pin")
}

func TestBotConfig_InsecureWithCAPath(t *testing.T) {
	cfg := &BotConfig{
		Insecure: true,
		Onboarding: OnboardingConfig{
			CAPath: "/tmp/invalid-path/some.crt",
		},
	}

	require.ErrorContains(t, cfg.CheckAndSetDefaults(), "ca-path")
}

func TestBotConfig_WithCAPathAndCAPins(t *testing.T) {
	cfg := &BotConfig{
		Insecure: false,
		Onboarding: OnboardingConfig{
			CAPath: "/tmp/invalid-path/some.crt",
			CAPins: []string{"123"},
		},
	}

	require.ErrorContains(t, cfg.CheckAndSetDefaults(), "mutually exclusive")
}

func TestBotConfig_ServicePartialCredentialLifetime(t *testing.T) {
	cfg := &BotConfig{
		Version:    V2,
		AuthServer: "example.teleport.sh:443",
		Services: []ServiceConfig{
			&IdentityOutput{
				CredentialLifetime: CredentialLifetime{TTL: 5 * time.Minute},
				Destination:        &DestinationMemory{},
			},
		},
	}
	require.ErrorContains(t, cfg.CheckAndSetDefaults(), "credential_ttl and renewal_interval")
}

func TestBotConfig_ServiceInvalidCredentialLifetime(t *testing.T) {
	cfg := &BotConfig{
		Version:    V2,
		AuthServer: "example.teleport.sh:443",
		Services: []ServiceConfig{
			&IdentityOutput{
				CredentialLifetime: CredentialLifetime{TTL: 5 * time.Minute},
				Destination:        &DestinationMemory{},
			},
		},
	}
	require.ErrorContains(t, cfg.CheckAndSetDefaults(), "credential_ttl and renewal_interval")
}

func TestBotConfig_DeprecatedCertificateTTL(t *testing.T) {
	t.Run("just deprecated option", func(t *testing.T) {
		const config = `
version: v2
certificate_ttl: 5m
`

		cfg, err := ReadConfig(strings.NewReader(config), false)
		require.NoError(t, err)
		require.Equal(t, 5*time.Minute, cfg.CredentialLifetime.TTL)
	})

	t.Run("both options", func(t *testing.T) {
		const config = `
version: v2
certificate_ttl: 5m
credential_ttl: 10m
`

		cfg, err := ReadConfig(strings.NewReader(config), false)
		require.NoError(t, err)
		require.Equal(t, 10*time.Minute, cfg.CredentialLifetime.TTL)
	})
}

func TestCredentialLifetimeValidate(t *testing.T) {
	testCases := map[string]struct {
		cfg     CredentialLifetime
		oneShot bool
		error   string
	}{
		"partial config": {
			cfg:   CredentialLifetime{TTL: 1 * time.Minute},
			error: "credential_ttl and renewal_interval must both be specified if either is",
		},
		"negative TTL": {
			cfg:   CredentialLifetime{TTL: -time.Minute, RenewalInterval: time.Minute},
			error: "credential_ttl must be positive",
		},
		"negative renewal interval": {
			cfg:   CredentialLifetime{TTL: time.Minute, RenewalInterval: -time.Minute},
			error: "renewal_interval must be positive",
		},
		"TTL less than renewal interval": {
			cfg:   CredentialLifetime{TTL: time.Minute, RenewalInterval: 2 * time.Minute},
			error: "TTL is shorter than the renewal interval",
		},
		"TTL less than renewal interval (one-shot)": {
			cfg:     CredentialLifetime{TTL: time.Minute, RenewalInterval: 2 * time.Minute},
			oneShot: true,
			error:   "",
		},
		"TTL too long": {
			cfg:   CredentialLifetime{TTL: defaults.MaxRenewableCertTTL * 2, RenewalInterval: time.Minute},
			error: "TTL exceeds the maximum TTL allowed",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.cfg.Validate(tc.oneShot)

			if tc.error == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.error)
			}
		})
	}
}

// TestBotConfig_Base64 ensures that config can be read from bas64 encoded YAML
func TestBotConfig_Base64(t *testing.T) {
	tests := []struct {
		name         string
		configBase64 string
		expected     BotConfig
	}{
		{
			name:         "minimal config, proxy server",
			configBase64: "dmVyc2lvbjogdjIKcHJveHlfc2VydmVyOiAiZXhhbXBsZS50ZWxlcG9ydC5zaDo0NDMiCm9uYm9hcmRpbmc6CiAgdG9rZW46ICJteS10b2tlbiIKICBqb2luX21ldGhvZDogInRva2VuIgpzZXJ2aWNlczoKLSB0eXBlOiBhcHBsaWNhdGlvbi10dW5uZWwKICBhcHBfbmFtZTogdGVzdGFwcAogIGxpc3RlbjogdGNwOi8vMTI3LjAuMC4xOjgwODA=",
			expected: BotConfig{
				Version:     V2,
				ProxyServer: "example.teleport.sh:443",
				Onboarding: OnboardingConfig{
					JoinMethod: "token",
					TokenValue: "my-token",
				},
				Services: []ServiceConfig{
					&ApplicationTunnelService{
						Listen:  "tcp://127.0.0.1:8080",
						AppName: "testapp",
					},
				},
			},
		},
		{
			name:         "minimal config, auth server",
			configBase64: "dmVyc2lvbjogdjIKYXV0aF9zZXJ2ZXI6ICJleGFtcGxlLnRlbGVwb3J0LnNoOjQ0MyIKb25ib2FyZGluZzoKICB0b2tlbjogIm15LXRva2VuIgogIGpvaW5fbWV0aG9kOiAidG9rZW4iCnNlcnZpY2VzOgotIHR5cGU6IGFwcGxpY2F0aW9uLXR1bm5lbAogIGFwcF9uYW1lOiB0ZXN0YXBwCiAgbGlzdGVuOiB0Y3A6Ly8xMjcuMC4wLjE6ODA4MA==",
			expected: BotConfig{
				Version:    V2,
				AuthServer: "example.teleport.sh:443",
				Onboarding: OnboardingConfig{
					JoinMethod: "token",
					TokenValue: "my-token",
				},
				Services: []ServiceConfig{
					&ApplicationTunnelService{
						Listen:  "tcp://127.0.0.1:8080",
						AppName: "testapp",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := ReadConfigFromBase64String(tt.configBase64, false)
			require.NoError(t, err)
			require.Equal(t, tt.expected, *cfg)
		})
	}
}

func TestBotConfig_NameValidation(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		cfg *BotConfig
		err string
	}{
		"duplicate names": {
			cfg: &BotConfig{
				Version: V2,
				Services: ServiceConfigs{
					&IdentityOutput{
						Name:        "foo",
						Destination: &DestinationMemory{},
					},
					&IdentityOutput{
						Name:        "foo",
						Destination: &DestinationMemory{},
					},
				},
			},
			err: `duplicate name: "foo`,
		},
		"reserved name": {
			cfg: &BotConfig{
				Version: V2,
				Services: ServiceConfigs{
					&IdentityOutput{
						Name:        "identity",
						Destination: &DestinationMemory{},
					},
				},
			},
			err: `service name "identity" is reserved for internal use`,
		},
		"invalid name": {
			cfg: &BotConfig{
				Version: V2,
				Services: ServiceConfigs{
					&IdentityOutput{
						Name:        "hello, world!",
						Destination: &DestinationMemory{},
					},
				},
			},
			err: `may only contain lowercase letters`,
		},
	}
	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			t.Parallel()
			require.ErrorContains(t, tc.cfg.CheckAndSetDefaults(), tc.err)
		})
	}
}
