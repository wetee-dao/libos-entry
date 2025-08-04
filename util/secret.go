package util

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"io"

	chain "github.com/wetee-dao/ink.go"
	"go.dedis.ch/kyber/v4"
	"go.dedis.ch/kyber/v4/suites"
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

// DecryptSecret decrypts a secret using the reader's secret key.
//
// Input:
//
//	ste                 - Crypto suite.
//	encScrt (rsG + K) - Encrypted key-slices.
//	dkgPk  (sG)         - Aggregate public key of DKG.
//	xncCmt (rsG + xsG)  - Re-encrypted schnorr-commit.
//	rdrSk  (x)          - Secret key of the reader.
//
// Output:
//
//	scrt - Recovered secret.
//	err - Error if decryption failed.
func DecryptSecret(
	ste suites.Suite,
	encScrt []kyber.Point,
	dkgPk kyber.Point,
	xncCmt kyber.Point,
	rdrSk kyber.Scalar,
) (
	scrt []byte,
	err error,
) {
	// To retrieve each key slice (Ki) from the encrypted key point (Ki + rsG),
	// we must deduct the encryption point (rsG). This can be inferred from
	// the re-encrypted schnorr-commit (rsG + xsG) by removing the product of
	// the reader's secret key (x) and the aggregate public key from the DKG (sG).
	xsG := ste.Point().Mul(rdrSk, dkgPk) // xsG = x * sG
	rsG := ste.Point().Sub(xncCmt, xsG)  // rsG = (rsG + xsG) - xsG

	for _, encKey := range encScrt {
		k := ste.Point().Sub(encKey, rsG) // K = (rsG + K) - rsG
		keyi, err := k.Data()
		if err != nil {
			return nil, fmt.Errorf("extract key share from key point: %w", err)
		}
		scrt = append(scrt, keyi...)
	}

	return scrt, err
}
