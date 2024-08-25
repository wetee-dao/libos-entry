package libos

import (
	"crypto/ed25519"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/wetee-dao/go-sdk/core"
	"github.com/wetee-dao/libos-entry/util"
)

func TestStartEntryServer(t *testing.T) {
	fs := &MockFs{}
	_, deployKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Error(err)
	}
	deploySinger, err := core.Ed25519PairFromPk(deployKey, 42)
	if err != nil {
		t.Error(err)
	}

	go startEntryServer(fs, &deploySinger, "")

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "https://0.0.0.0:8888/report", nil)
	time.Sleep(time.Second * 2)
	http.DefaultServeMux.ServeHTTP(w, req)
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

func (l *MockFs) VerifyReport(workerReport *util.TeeParam) (*util.TeeReport, error) {
	return nil, nil
}

func (l *MockFs) IssueReport(pk *core.Signer, data []byte) (*util.TeeParam, error) {
	return &util.TeeParam{}, nil
}

func (l *MockFs) SetPassword(password string) {

}
