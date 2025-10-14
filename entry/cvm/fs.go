package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/afero"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/libos-entry/model"
)

type Fs struct {
	afero.OsFs
}

// Read implements util.Fs.
func (e *Fs) ReadFile(filename string) ([]byte, error) {
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
func (e *Fs) WriteFile(filename string, data []byte, perm os.FileMode) error {
	// 加密数据
	val, err := e.Encrypt(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt Key: %v", err)
	}

	return afero.WriteFile(e, filename, val, perm)
}

// Decrypt implements libos.SecretFunction.
func (e *Fs) Decrypt(val []byte) ([]byte, error) {
	return val, nil
}

// Encrypt implements libos.SecretFunction.
func (e *Fs) Encrypt(val []byte) ([]byte, error) {
	return val, nil
}

// IssueReport implements libos.TeeFunction.
func (i *Fs) IssueReport(pk chain.Signer, call *model.TeeCall) error {
	return model.IssueReport(&pk, call)
}

// VerifyReport implements libos.TeeFunction.
func (e *Fs) VerifyReport(reportData *model.TeeCall) (*model.TeeVerifyResult, error) {
	// 检查时间戳，超过 30s 签名过期
	if reportData.Time+30 < time.Now().Unix() {
		// return nil, errors.New("report expired")
	}

	return model.VerifyReport(reportData)
}
