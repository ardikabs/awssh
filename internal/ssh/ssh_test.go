package ssh_test

import (
	"awssh/config"
	"awssh/internal/logging"
	. "awssh/internal/ssh"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh/agent"
)

type mockSSHAgent struct {
	agent.ExtendedAgent
}

func (m mockSSHAgent) List() ([]*agent.Key, error) {
	return []*agent.Key{}, nil
}

func (m mockSSHAgent) Add(key agent.AddedKey) error {
	return nil
}

func TestMain(m *testing.M) {
	config.Load()
	logging.NewLogger(false)
	code := m.Run()
	os.Exit(code)
}

func TestNewSession(t *testing.T) {
	sess, err := NewSession(mockSSHAgent{}, "i-123456789abc")
	assert.Nil(t, err)
	assert.NotNil(t, sess)
}
