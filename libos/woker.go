package libos

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// Worker 请求通道
type WorkerChannel struct {
	TlsConfig *tls.Config
}

func (w *WorkerChannel) Get(url string) ([]byte, error) {
	client := http.Client{Transport: &http.Transport{TLSClientConfig: w.TlsConfig}}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (w *WorkerChannel) Post(url string, json string) ([]byte, error) {
	client := http.Client{Transport: &http.Transport{TLSClientConfig: w.TlsConfig}}
	payload := strings.NewReader(json)
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bt, err := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(bt))
	}

	return bt, err
}
