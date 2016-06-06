package keygen

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

type keyPair struct {
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

func (k keyPair) String() string {
	return fmt.Sprintf("Private key: %v\nPublic key: %v", k.Private, k.Public)
}

// NewKeyPair creates a new pair of RSA keys.
// For keys use 2048 bit keys.
func NewKeyPair() (*keyPair, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	privKey.Precompute()

	pubKey := &privKey.PublicKey
	return &keyPair{Private: privKey, Public: pubKey}, nil
}
