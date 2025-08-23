package model

import (
	"encoding/binary"
	"errors"
	"os"
	"time"

	"github.com/edgelesssys/ego/enclave"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
)

var (
	TeeType            = 9999
	IsEgo              = false
	SevGuestDevicePath = "/dev/sev-guest"
	TdxGuestDevicePath = "/dev/tdx_guest"
)

func init() {
	if CheckExists(SevGuestDevicePath) {
		util.LogWithGreen("TEE TYPE", "SEV-SNP")
		TeeType = 1
		return
	}

	if CheckExists(TdxGuestDevicePath) {
		util.LogWithGreen("TEE TYPE", "TDX")
		TeeType = 2
		return
	}

	if _, err := enclave.GetSelfReport(); err == nil {
		util.LogWithGreen("TEE TYPE", "EGO SGX")
		TeeType = 0
		IsEgo = true
		return
	}

	util.LogWithGreen("TEE TYPE", "NO TEE")
}

// issue report
func IssueReport(pk *chain.Signer, call *TeeCall) error {
	switch TeeType {
	case 0:
		return SgxIssue(pk, call)
	case 1:
		return SnpIssue(pk, call)
	default:
		timestamp := time.Now().Unix()
		call.Time = timestamp
		call.TeeType = 9999
		call.Caller = pk.PublicKey
		return nil
	}
}

// verify report
func VerifyReport(reportData *TeeCall) (*TeeVerifyResult, error) {
	switch reportData.TeeType {
	case 0:
		if IsEgo {
			return SgxVerify(reportData)
		} else {
			return ClientSgxVerify(reportData)
		}
	case 1:
		return SnpVerify(reportData)
	}

	return nil, errors.New("unknown tee type")
}

// int to bytes
func Int64ToBytes(time int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(time))
	return b
}

// bytes to int
func BytesToInt64(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}

// Exists 检查路径是否存在（无论是文件还是目录）
func CheckExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}
