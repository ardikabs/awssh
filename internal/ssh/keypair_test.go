package ssh_test

import (
	"testing"

	. "awssh/internal/ssh"

	"github.com/stretchr/testify/assert"
)

func TestNewKeyPair(t *testing.T) {

	t.Run("error when keysize is not factor 2", func(t *testing.T) {
		keypair, err := NewKeyPair(5)

		assert.NotNil(t, err)
		assert.Nil(t, keypair)
	})

	t.Run("successfully create keypair when keysize is factor 2, e.g 1024/2048/3072/4096", func(t *testing.T) {
		keypair, err := NewKeyPair(1024)

		assert.Nil(t, err)
		assert.NotNil(t, keypair)
	})
}
