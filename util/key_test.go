package util

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	chain "github.com/wetee-dao/go-sdk"
)

func TestSetKey(t *testing.T) {
	// Create a mock filesystem
	mockFs := &MockFs{}

	// Call the SetKey function
	key, err := GetKey(mockFs, "./xx")
	if err != nil {
		t.Errorf("SetKey returned an error: %s", err)
	}
	// Call the SetKey function
	key2, err := GetKey(mockFs, "./xx")
	if err != nil {
		t.Errorf("SetKey returned an error: %s", err)
	}

	if key.SS58Address(42) != key2.SS58Address(42) {
		t.Errorf("SetKey returned different keys")
	}

	os.Remove("./xx")
}

type MockFs struct {
	afero.OsFs
}

// Read implements util.Fs.
func (f *MockFs) ReadFile(filename string) ([]byte, error) {
	return afero.ReadFile(f, filename)
}

// Write implements util.Fs.
func (f *MockFs) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return afero.WriteFile(f, filename, data, perm)
}

func (l *MockFs) VerifyReport(reportData *TeeParam) (*TeeReport, error) {
	return nil, nil
}

func (l *MockFs) IssueReport(pk *chain.Signer, data []byte) (*TeeParam, error) {
	return &TeeParam{}, nil
}

func (l *MockFs) SetPassword(password string) {

}
