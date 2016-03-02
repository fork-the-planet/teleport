package teleagent

import (
	"io"
	"net"
	"time"

	"github.com/gravitational/teleport/lib/auth/native"
	"github.com/gravitational/teleport/lib/utils"
	"github.com/gravitational/teleport/lib/web"

	log "github.com/Sirupsen/logrus"
	"github.com/gravitational/trace"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type TeleAgent struct {
	agent    agent.Agent
	insecure bool
}

func NewTeleAgent(insecure bool) *TeleAgent {
	ta := TeleAgent{
		agent:    agent.NewKeyring(),
		insecure: insecure,
	}
	return &ta
}

func (a *TeleAgent) Start(agentAddr string) error {
	addr, err := utils.ParseAddr(agentAddr)
	if err != nil {
		return trace.Wrap(err)
	}

	l, err := net.Listen(addr.AddrNetwork, addr.Addr)
	if err != nil {
		return trace.Wrap(err)
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Errorf(err.Error())
				continue
			}
			go func() {
				if err := agent.ServeAgent(a.agent, conn); err != nil {
					if err != io.EOF {
						log.Errorf(err.Error())
					}
				}
			}()
		}
	}()

	return nil
}

func (a *TeleAgent) Login(proxyAddr string,
	user string, pass string, hotpToken string,
	ttl time.Duration) error {

	priv, pub, err := native.New().GenerateKeyPair("")
	if err != nil {
		return trace.Wrap(err)
	}

	login, err := web.SSHAgentLogin(proxyAddr, user, pass, hotpToken,
		pub, ttl, a.insecure)
	if err != nil {
		return trace.Wrap(err)
	}

	pcert, _, _, _, err := ssh.ParseAuthorizedKey(login.Cert)
	if err != nil {
		return trace.Wrap(err)
	}

	pk, err := ssh.ParseRawPrivateKey(priv)
	if err != nil {
		return trace.Wrap(err)
	}
	addedKey := agent.AddedKey{
		PrivateKey:       pk,
		Certificate:      pcert.(*ssh.Certificate),
		Comment:          "",
		LifetimeSecs:     0,
		ConfirmBeforeUse: false,
	}
	if err := a.agent.Add(addedKey); err != nil {
		return trace.Wrap(err)
	}

	return nil
}

func (a *TeleAgent) AuthMethod() ssh.AuthMethod {
	return ssh.PublicKeysCallback(a.agent.Signers)
}

const (
	DefaultAgentAddress = "unix:///tmp/teleport.agent.sock"
)
