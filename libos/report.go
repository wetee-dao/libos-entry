package libos

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"github.com/edgelesssys/ego/attestation"
	"github.com/pkg/errors"
	"github.com/vedhavyas/go-subkey/v2"
	"github.com/wetee-dao/libos-entry/util"
)

// CreateCertificate create certificate
func CreateCertificate(appId string) ([]byte, crypto.PrivateKey) {
	template := &x509.Certificate{
		SerialNumber: &big.Int{},
		Subject:      pkix.Name{CommonName: "wetee.app"},
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		DNSNames:     []string{appId},
	}
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	cert, _ := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	return cert, priv
}

// VerifyReport verify report
func VerifyReport(workerReport *util.TeeParam, fs util.Fs) (*attestation.Report, error) {
	// decode address
	_, signer, err := subkey.SS58Decode(workerReport.Address)
	if err != nil {
		return nil, errors.Wrap(err, "SS58 decode")
	}

	return fs.VerifyReport(workerReport.Report, workerReport.Data, signer, workerReport.Time)
}
