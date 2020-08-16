package ssh

import (
	"fmt"
	"net"
	"os"

	"awssh/config"

	"golang.org/x/crypto/ssh/agent"
)

// Session is
type Session struct {
	PublicKey string
}

// NewSession is TODO:
func NewSession(instanceID string) (session *Session, err error) {
	sshSocket := os.Getenv("SSH_AUTH_SOCK")
	conn, err := net.Dial("unix", sshSocket)

	if err != nil {
		return nil, fmt.Errorf("Failed to establish a connection to SSH_AUTH_SOCK: (%v)", err)
	}

	agentClient := agent.NewClient(conn)

	existKeys, err := agentClient.List()

	if err != nil {
		return nil, err
	}

	var publicKey string

	if len(existKeys) == 0 {
		keypair, err := NewKeyPair(2048)

		appConfig := config.Get()

		tmpKey := agent.AddedKey{
			PrivateKey:       keypair.PrivateKey,
			Comment:          fmt.Sprintf("awssh-temporary-ssh-keypair:%s:%s", appConfig.SSHUsername, instanceID),
			LifetimeSecs:     30,
			ConfirmBeforeUse: false,
		}

		err = agentClient.Add(tmpKey)
		if err != nil {
			return nil, fmt.Errorf("Unable to add ssh keypair to ssh agent: (%v)", err)

		}

		publicKey = string(keypair.PublicKey)
	} else {
		// use the first ssh-key registered in ssh-agent to be loaded
		publicKey = fmt.Sprintf("%s", existKeys[0])
	}

	return &Session{
		PublicKey: publicKey,
	}, nil
}
