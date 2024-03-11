package libos

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/spf13/afero"
)

func TestStartEntryServer(t *testing.T) {
	fs := &MockFs{}
	// 获取本地证书
	// Get local certificate
	certBytes, priv, report, err := GetLocalReport("xxx", fs)
	if err != nil {
		t.Errorf("GetRemoteReport: %v", err)
	}

	go startEntryServer(certBytes, priv, report)

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

func (l *MockFs) VerifyReport(reportBytes, certBytes, signer []byte) error {
	return nil
}

func (l *MockFs) IssueReport(data []byte) ([]byte, error) {
	return data, nil
}

func (l *MockFs) SetPassword(password string) {

}
