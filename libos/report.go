package libos

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"time"

	"github.com/edgelesssys/ego/attestation"
	"github.com/pkg/errors"
	"github.com/wetee-dao/libos-entry/util"
)

// GetRemoteReport get remote report
func GetLocalReport(appId string, fs util.Fs) ([]byte, crypto.PrivateKey, []byte, error) {
	cert, priv := CreateCertificate(appId)
	hash := sha256.Sum256(cert)
	report, err := fs.IssueReport(hash[:])
	return cert, priv, report, err
}

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

func VerifyReport(workerReportWrap []byte, fs util.Fs) (*attestation.Report, error) {
	workerReport := map[string]string{}
	err := json.Unmarshal(workerReportWrap, &workerReport)
	if err != nil {
		return nil, errors.Wrap(err, "Unmarshal worker report")
	}

	report, err := hex.DecodeString(workerReport["report"])
	if err != nil {
		return nil, errors.Wrap(err, "Hex decode worker report")
	}

	return fs.VerifyReport(report, nil, nil)
}
