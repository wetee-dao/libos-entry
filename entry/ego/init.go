package ego

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/edgelesssys/ego/ecrypto"
	"github.com/spf13/afero"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/libos-entry/libos"
	"github.com/wetee-dao/libos-entry/model"
)

// var Fs *EgoFs

func InitEgo() error {
	hostfs := &EgoFs{}
	return libos.PreLoad(hostfs, false)
}

type EgoFs struct {
	afero.OsFs
	report     []byte
	lastReport int64
}

// Read implements util.Fs.
func (e *EgoFs) ReadFile(filename string) ([]byte, error) {
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
func (e *EgoFs) WriteFile(filename string, data []byte, perm os.FileMode) error {

	// 加密数据
	val, err := e.Encrypt(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt Key: %v", err)
	}

	return afero.WriteFile(e, filename, val, perm)
}

// Decrypt implements libos.SecretFunction.
func (e *EgoFs) Decrypt(val []byte) ([]byte, error) {
	var additionalData []byte = nil
	return ecrypto.Unseal(val, additionalData)
}

// Encrypt implements libos.SecretFunction.
func (e *EgoFs) Encrypt(val []byte) ([]byte, error) {
	var additionalData []byte = nil
	return ecrypto.SealWithProductKey(val, additionalData)
}

// IssueReport implements libos.TeeFunction.
func (i *EgoFs) IssueReport(pk chain.Signer, call *model.TeeCall) error {
	return model.IssueReport(&pk, call)
}

// VerifyReport implements libos.TeeFunction.
func (e *EgoFs) VerifyReport(workerReport *model.TeeCall) (*model.TeeVerifyResult, error) {
	// 检查时间戳，超过 30s 签名过期
	if workerReport.Time+30 < time.Now().Unix() {
		return nil, errors.New("report expired")
	}

	return model.VerifyReport(workerReport)
}
