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

func NewKeyPair() (*keyPair, error) {
	priv_key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	priv_key.Precompute()

	pub_key := &priv_key.PublicKey
	return &keyPair{Private: priv_key, Public: pub_key}, nil
}
