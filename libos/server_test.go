package libos

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/edgelesssys/ego/attestation"
	"github.com/spf13/afero"
	"github.com/wetee-dao/go-sdk/core"
)

func TestStartEntryServer(t *testing.T) {
	fs := &MockFs{}

	go startEntryServer(fs)

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

func (l *MockFs) VerifyReport(reportBytes, certBytes, signer []byte, t int64) (*attestation.Report, error) {
	return nil, nil
}

func (l *MockFs) IssueReport(pk *core.Signer, data []byte) ([]byte, int64, error) {
	return data, 0, nil
}

func (l *MockFs) SetPassword(password string) {

}
