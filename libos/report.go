package libos

import (
	"bytes"
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
	// decode time
	timestamp := workerReport.Time

	// 检查时间戳，超过 30s 签名过期
	if timestamp+30 < time.Now().Unix() {
		return nil, errors.New("Report expired")
	}

	// decode report
	report := workerReport.Report

	// decode address
	_, signer, err := subkey.SS58Decode(workerReport.Address)
	if err != nil {
		return nil, errors.Wrap(err, "SS58 decode")
	}

	// 构建验证数据
	var buf bytes.Buffer
	buf.Write(util.Int64ToBytes(timestamp))
	buf.Write(signer)

	return fs.VerifyReport(report, buf.Bytes(), signer)
}
