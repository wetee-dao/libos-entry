package model

import (
	"encoding/binary"
	"errors"
	"os"
	"time"

	chain "github.com/wetee-dao/ink.go"
)

var (
	SevGuestDevicePath = "/dev/sev-guest"
	TdxGuestDevicePath = "/dev/tdx_guest"
)

func IssueReport(pk *chain.Signer, call *TeeCall) error {
	if CheckExists(SevGuestDevicePath) {
		return SnpIssue(pk, call)
	}

	if CheckExists("/dev/sgx") {
		return SgxIssue(pk, call)
	}

	timestamp := time.Now().Unix()
	call.Time = timestamp
	call.TeeType = 9999
	call.Caller = pk.PublicKey
	return nil
}

func VerifyReport(reportData *TeeCall) (*TeeVerifyResult, error) {
	switch reportData.TeeType {
	case 0:
		if CheckExists("/dev/sgx") {
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
