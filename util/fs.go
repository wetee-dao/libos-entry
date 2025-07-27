package util

import (
	"os"

	"github.com/spf13/afero"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/libos-entry/model"
)

type Fs interface {
	afero.Fs
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error

	// VerifyReport(reportBytes, data, signer []byte, t int64) (*attestation.Report, error)
	VerifyReport(workerReport *model.TeeCall) (*model.TeeVerifyResult, error)
	IssueReport(pk chain.Signer, data *model.TeeCall) error
}
