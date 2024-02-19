package libos

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"github.com/edgelesssys/ego/enclave"
)

func GetRemoteReport() ([]byte, crypto.PrivateKey, []byte, error) {
	cert, priv := CreateCertificate()
	hash := sha256.Sum256(cert)
	report, err := enclave.GetRemoteReport(hash[:])
	return cert, priv, report, err
}

func CreateCertificate() ([]byte, crypto.PrivateKey) {
	template := &x509.Certificate{
		SerialNumber: &big.Int{},
		Subject:      pkix.Name{CommonName: "wetee.app"},
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		DNSNames:     []string{"localhost"},
	}
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	cert, _ := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	return cert, priv
}
