package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/afero"
	chain "github.com/wetee-dao/ink.go"
	inkutil "github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/libos-entry/libos"
	"github.com/wetee-dao/libos-entry/model"
)

type CvmServer struct {
}

func init() {
	TEEServerImpl = CvmServer{}
}

func (CvmServer) start(req *CrossRequest) CrossResponse {
	envs := map[string]string{}
	err := json.Unmarshal(req.env, &envs)
	if err != nil {
		return CrossResponse{code: 1, data: []byte(err.Error())}
	}

	hostfs := &Fs{}
	err = libos.PreLoadFromInitData(hostfs, envs, false)
	if err != nil {
		inkutil.LogWithGray("PreLoadFromInitData", err.Error())
		return CrossResponse{code: 1, data: []byte(err.Error())}
	}

	return CrossResponse{code: 0}
}

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
	timestamp := time.Now().Unix()
	call.Time = timestamp
	call.TeeType = 1

	call.Caller = pk.PublicKey
	return nil
}

// VerifyReport implements libos.TeeFunction.
func (e *Fs) VerifyReport(workerReport *model.TeeCall) (*model.TeeVerifyResult, error) {
	// 检查时间戳，超过 30s 签名过期
	if workerReport.Time+30 < time.Now().Unix() {
		return nil, errors.New("report expired")
	}

	// return model.VerifyReport(workerReport)
	return &model.TeeVerifyResult{
		TeeType:       workerReport.TeeType,
		CodeSigner:    []byte{},
		CodeSignature: []byte{},
		CodeProductId: []byte{},
	}, nil
}
