package ssh

import (
	"fmt"
	"net"
	"os"

	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	"awssh/config"
	"awssh/internal/logging"
)

// Session represent an SSH data model consist of a SSH PublicKey
type Session struct {
	PublicKey string
}

// NewSession creates a new SSH session from instanceID
// This method will determine to select whether need to create a new temporary ssh keypair
// or used the first existing key given from ssh-agent
func NewSession(instanceID string) (session *Session, err error) {
	appConfig := config.Get()

	logger := logging.Get()

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
		if err != nil {
			return nil, err
		}

		tmpSSHKeyPair := agent.AddedKey{
			PrivateKey:       keypair.PrivateKey,
			Comment:          fmt.Sprintf("awssh-temporary-ssh-keypair:%s:%s", appConfig.SSHUsername, instanceID),
			LifetimeSecs:     30,
			ConfirmBeforeUse: false,
		}

		err = agentClient.Add(tmpSSHKeyPair)
		if err != nil {
			return nil, fmt.Errorf("Unable to add ssh keypair to ssh agent: (%v)", err)
		}

		logger.Debugf("Create temporary ssh-rsa keypair (%s)", gossh.FingerprintSHA256(keypair.PublicKey))
		publicKeySerialized := gossh.MarshalAuthorizedKey(keypair.PublicKey)
		publicKey = string(publicKeySerialized)
	} else {
		logger.Debugf("Use existing ssh-rsa keypair from ssh-agent (%s)", gossh.FingerprintSHA256(existKeys[0]))
		publicKey = fmt.Sprint(existKeys[0])
	}

	return &Session{
		PublicKey: publicKey,
	}, nil
}
