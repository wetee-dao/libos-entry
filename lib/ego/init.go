package ego

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"

	"github.com/edgelesssys/ego/ecrypto"
	"github.com/edgelesssys/ego/enclave"
)

func InitEgo(chainAddr string) error {
	// hostfs := afero.NewOsFs()
	// return libos.PreLoad(chainAddr, hostfs, &EgoSf{})
	return nil
}

type EgoSf struct {
}

// Decrypt implements libos.SecretFunction.
func (l *EgoSf) Decrypt(val []byte) ([]byte, error) {
	return ecrypto.Unseal(val, nil)
}

// Encrypt implements libos.SecretFunction.
func (l *EgoSf) Encrypt(val []byte) ([]byte, error) {
	return ecrypto.SealWithProductKey(val, nil)
}

func (l *EgoSf) VerifyReport(reportBytes, certBytes, signer []byte) error {
	report, err := enclave.VerifyRemoteReport(reportBytes)
	if err != nil {
		return err
	}

	hash := sha256.Sum256(certBytes)
	if !bytes.Equal(report.Data[:len(hash)], hash[:]) {
		return errors.New("report data does not match the certificate's hash")
	}

	// You can either verify the UniqueID or the tuple (SignerID, ProductID, SecurityVersion, Debug).

	if report.SecurityVersion < 2 {
		return errors.New("invalid security version")
	}
	if binary.LittleEndian.Uint16(report.ProductID) != 1234 {
		return errors.New("invalid product")
	}
	if !bytes.Equal(report.SignerID, signer) {
		return errors.New("invalid signer")
	}

	// For production, you must also verify that report.Debug == false

	return nil
}
