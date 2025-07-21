package util

import (
	"crypto/ed25519"
	"crypto/rand"
	"io"

	chain "github.com/wetee-dao/ink.go"
)

// 去中心化的机密注入
type Secrets struct {
	Envs  map[string]string
	Files map[string][]byte
}

func GenerateKeyPair(src io.Reader) (*chain.Signer, ed25519.PrivateKey, ed25519.PublicKey, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, nil, err
	}

	kr, err := chain.Ed25519PairFromPk(privKey, 42)
	if err != nil {
		return nil, nil, nil, err
	}

	return &kr, privKey, pubKey, nil
}
