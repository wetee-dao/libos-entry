package util

import (
	"os"

	"github.com/edgelesssys/ego/attestation"
	"github.com/spf13/afero"
	"github.com/wetee-dao/go-sdk/core"
)

type Fs interface {
	afero.Fs
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error

	VerifyReport(reportBytes, data, signer []byte) (*attestation.Report, error)
	IssueReport(pk *core.Signer, data []byte) ([]byte, int64, error)
}
