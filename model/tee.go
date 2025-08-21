package model

import (
	"encoding/binary"
	"errors"
	"os"

	chain "github.com/wetee-dao/ink.go"
)

func IssueReport(pk *chain.Signer, call *TeeCall, teeType uint32) error {
	switch teeType {
	case 0:
		return SgxIssue(pk, call)
	case 1:
		return SnpIssue(pk, call)
	}

	return errors.New("unknown tee type")
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
