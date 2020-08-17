package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	gossh "golang.org/x/crypto/ssh"
)

// KeyPair represent of a SSH-RSA KeyPair data
type KeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  gossh.PublicKey
}

// NewKeyPair creates a new KeyPair from key size
func NewKeyPair(keysize int) (keypair *KeyPair, err error) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, keysize)

	if err != nil {
		return nil, fmt.Errorf("Failed to generate RSA key: (%v)", err)
	}

	publicKey, err := gossh.NewPublicKey(&rsaKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("Failed to create ssh public key: (%v)", err)
	}

	return &KeyPair{
		PrivateKey: rsaKey,
		PublicKey:  publicKey,
	}, nil
}
