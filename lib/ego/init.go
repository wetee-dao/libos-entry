package ego

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"github.com/edgelesssys/ego/ecrypto"
	"github.com/edgelesssys/ego/enclave"
	"github.com/spf13/afero"
	"github.com/wetee-dao/libos-entry/libos"
)

func InitEgo(chainAddr string) error {
	hostfs := &EgoSf{}
	return libos.PreLoad(chainAddr, hostfs)
}

type EgoSf struct {
	afero.OsFs
	password string
}

// Read implements util.Fs.
func (e *EgoSf) ReadFile(filename string) ([]byte, error) {
	bt, err := afero.ReadFile(e, filename)
	if err != nil {
		return nil, err
	}

	// 解密数据
	keyBytes, err := e.Decrypt(bt)
	if err != nil {
		return nil, err
	}

	return keyBytes, nil
}

// Write implements util.Fs.
func (e *EgoSf) WriteFile(filename string, data []byte, perm os.FileMode) error {

	// 加密数据
	val, err := e.Encrypt(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt Key: %v", err)
	}

	return afero.WriteFile(e, filename, val, perm)
}

// Decrypt implements libos.SecretFunction.
func (e *EgoSf) Decrypt(val []byte) ([]byte, error) {
	var additionalData []byte = nil
	if len(e.password) != 0 {
		additionalData = []byte(e.password)
	}
	return ecrypto.Unseal(val, additionalData)
}

// Encrypt implements libos.SecretFunction.
func (e *EgoSf) Encrypt(val []byte) ([]byte, error) {
	var additionalData []byte = nil
	if len(e.password) != 0 {
		additionalData = []byte(e.password)
	}
	return ecrypto.SealWithProductKey(val, additionalData)
}

func (e *EgoSf) VerifyReport(reportBytes, certBytes, signer []byte) error {
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
