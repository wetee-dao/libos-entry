package util

import (
	"os"

	"github.com/edgelesssys/ego/attestation"
	"github.com/spf13/afero"
)

type Fs interface {
	afero.Fs
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error

	SetPassword(password string)

	VerifyReport(reportBytes, certBytes, signer []byte) (*attestation.Report, error)
	IssueReport(cert []byte) ([]byte, error)
}
