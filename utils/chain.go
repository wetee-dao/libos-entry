package utils

import (
	"encoding/hex"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/vedhavyas/go-subkey/v2/sr25519"
)

// 获取签名密钥
// GetKey get signer key
func GetSignerKey() (*signature.KeyringPair, error) {
	kr, err := sr25519.Scheme{}.Generate()
	if err != nil {
		return nil, err
	}

	uri := hex.EncodeToString(kr.Seed())

	return &signature.KeyringPair{
		URI:       uri,
		Address:   kr.SS58Address(42),
		PublicKey: kr.Public(),
	}, nil
}
