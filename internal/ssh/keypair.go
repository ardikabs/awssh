package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	gossh "golang.org/x/crypto/ssh"
)

// KeyPair is a
type KeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  []byte
}

// NewKeyPair is a TODO:
func NewKeyPair(keysize int) (keypair *KeyPair, err error) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, keysize)

	if err != nil {
		return nil, fmt.Errorf("Failed to generate RSA key: (%v)", err)
	}

	publicKey, err := gossh.NewPublicKey(&rsaKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("Failed to create ssh public key: (%v)", err)
	}

	publicKeySerialized := gossh.MarshalAuthorizedKey(publicKey)

	return &KeyPair{
		PrivateKey: rsaKey,
		PublicKey:  publicKeySerialized,
	}, nil
}
