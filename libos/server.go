package libos

import (
	"crypto"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// 创建一个专门用于为外接用于证明和获取证明的服务
func startEntryServer(cert []byte, priv crypto.PrivateKey, report []byte) error {
	tlsCfg := tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{cert},
				PrivateKey:  priv,
			},
		},
	}

	router := chi.NewRouter()
	router.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]string{
			"report": hex.EncodeToString(report),
		}
		bt, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(bt)
	})

	router.HandleFunc("/set_password", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]string{
			"report": hex.EncodeToString(report),
		}
		bt, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(bt)
	})

	server := &http.Server{Addr: ":8888", Handler: router, TLSConfig: &tlsCfg}
	fmt.Println("Start entry secret listening https://0.0.0.0:8888 ...")
	err := server.ListenAndServeTLS("", "")
	return err
}
