package util

import (
	"crypto/ed25519"
	"io"

	libp2pCrypto "github.com/libp2p/go-libp2p/core/crypto"
	chain "github.com/wetee-dao/go-sdk"
)

// 去中心化的机密注入
type Secrets struct {
	Envs  map[string]string
	Files map[string][]byte
}

func GenerateKeyPair(src io.Reader) (*chain.Signer, error) {
	sk, _, err := libp2pCrypto.GenerateKeyPairWithReader(libp2pCrypto.Ed25519, 0, src)
	if err != nil {
		return nil, err
	}

	bt, err := sk.Raw()
	if err != nil {
		return nil, err
	}

	var ed25519Key ed25519.PrivateKey = bt
	kr, err := chain.Ed25519PairFromPk(ed25519Key, 42)
	if err != nil {
		return nil, err
	}

	return &kr, nil
}
