package ssh

import (
	"crypto/rand"
	"crypto/rsa"

	gossh "golang.org/x/crypto/ssh"
)

// KeyPair is a
type KeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  []byte
}

// NewKeyPair is a
func NewKeyPair(keysize int) (keypair *KeyPair, err error) {
	// TODO: aws.NewKeyPair proper docs
	rsaKey, err := rsa.GenerateKey(rand.Reader, keysize)

	if err != nil {
		return nil, err
	}

	publicKey, err := gossh.NewPublicKey(&rsaKey.PublicKey)
	if err != nil {
		return nil, err
	}

	publicKeySerialized := gossh.MarshalAuthorizedKey(publicKey)

	return &KeyPair{
		PrivateKey: rsaKey,
		PublicKey:  publicKeySerialized,
	}, nil
}
