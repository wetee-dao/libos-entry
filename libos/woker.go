package libos

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

// Worker 请求通道
type WorkerChannel struct {
	TlsConfig *tls.Config
}

func (w *WorkerChannel) Get(url string) ([]byte, error) {
	client := http.Client{Transport: &http.Transport{TLSClientConfig: w.TlsConfig}, Timeout: time.Minute}
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
	client := http.Client{Transport: &http.Transport{TLSClientConfig: w.TlsConfig}, Timeout: time.Minute}
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

func (w *WorkerChannel) PostBt(url string, data []byte) ([]byte, error) {
	client := http.Client{Transport: &http.Transport{TLSClientConfig: w.TlsConfig}, Timeout: time.Minute}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/x-protobuf")
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
