package util

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

// Ed25519Cert 生成 Ed25519 自签名证书
func Ed25519Cert(
	commonName string,
	ips []net.IP,
	dns []string,
	privKey ed25519.PrivateKey,
	pubKey ed25519.PublicKey,
) (cert, key []byte, der []byte, err error) {
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IPAddresses:           ips,
		DNSNames:              dns,
	}

	// Create certificate
	der, err = x509.CreateCertificate(rand.Reader, &template, &template, pubKey, privKey)
	if err != nil {
		return nil, nil, nil, err
	}

	// 证书 PEM
	cert = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})

	// 私钥 PEM（PKCS#8格式）
	keyBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return nil, nil, nil, err
	}
	key = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyBytes})

	return cert, key, der, nil
}
