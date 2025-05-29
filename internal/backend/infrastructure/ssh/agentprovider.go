package ssh

import (
	"fmt"
	"os"

	"github.com/nightnoryu/bkit/internal/backend/app/ssh"
)

const (
	sshAuthSock = "SSH_AUTH_SOCK"
)

func NewAgentProvider() (ssh.AgentProvider, error) {
	socket, found := os.LookupEnv(sshAuthSock)
	if !found {
		return nil, fmt.Errorf("ssh auth socket via env %s not found", sshAuthSock)
	}

	return &agentProvider{defaultAgent: socket}, nil
}

type agentProvider struct {
	defaultAgent string
}

func (provider agentProvider) Default() string {
	return provider.defaultAgent
}
