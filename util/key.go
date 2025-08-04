package util

import (
	"crypto/ed25519"
	"crypto/sha512"

	"go.dedis.ch/kyber/v4"
	"go.dedis.ch/kyber/v4/suites"
)

func Ed25519Scalar(suite suites.Suite, buf ed25519.PrivateKey) kyber.Scalar {
	// hash seed and clamp bytes
	digest := sha512.Sum512(buf[:32])
	digest[0] &= 0xf8
	digest[31] &= 0x7f
	digest[31] |= 0x40
	return suite.Scalar().SetBytes(digest[:32])
}
