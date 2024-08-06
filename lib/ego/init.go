package ego

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/edgelesssys/ego/attestation"
	"github.com/edgelesssys/ego/attestation/tcbstatus"
	"github.com/edgelesssys/ego/ecrypto"
	"github.com/edgelesssys/ego/enclave"
	"github.com/spf13/afero"
	"github.com/vedhavyas/go-subkey/v2/ed25519"
	"github.com/wetee-dao/go-sdk/core"
	"github.com/wetee-dao/libos-entry/libos"
	"github.com/wetee-dao/libos-entry/util"
)

func InitEgo() error {
	hostfs := &EgoFs{}
	return libos.PreLoad(hostfs)
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

func (i *EgoFs) IssueReport(pk *core.Signer, data []byte) ([]byte, int64, error) {
	timestamp := time.Now().Unix()
	if i.report != nil && i.lastReport+30 > timestamp {
		return i.report, i.lastReport, nil
	}

	var buf bytes.Buffer
	buf.Write(util.Int64ToBytes(timestamp))
	buf.Write(pk.PublicKey)
	if len(data) > 0 {
		buf.Write(data)
	}
	sig, err := pk.Sign(buf.Bytes())
	if err != nil {
		return nil, 0, err
	}

	report, err := enclave.GetRemoteReport(sig)
	if err != nil {
		return nil, 0, err
	}

	i.lastReport = timestamp
	i.report = report
	return report, timestamp, nil
}

func (e *EgoFs) VerifyReport(reportBytes, msgBytes, signer []byte, timestamp int64) (*attestation.Report, error) {
	// 检查时间戳，超过 30s 签名过期
	if timestamp+30 < time.Now().Unix() {
		return nil, errors.New("report expired")
	}

	report, err := enclave.VerifyRemoteReport(reportBytes)
	if err == attestation.ErrTCBLevelInvalid {
		fmt.Printf("Warning: TCB level is invalid: %v\n%v\n", report.TCBStatus, tcbstatus.Explain(report.TCBStatus))
	} else if err != nil {
		return nil, err
	}

	pubkey, err := ed25519.Scheme{}.FromPublicKey(signer)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.Write(util.Int64ToBytes(timestamp))
	buf.Write(signer)
	if len(msgBytes) > 0 {
		buf.Write(msgBytes)
	}

	sig := report.Data

	if !pubkey.Verify(buf.Bytes(), sig) {
		return nil, errors.New("invalid sgx report")
	}

	if report.Debug {
		return nil, errors.New("debug mode is not allowed")
	}

	return &report, nil
}
