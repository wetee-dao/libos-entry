package libos

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/wetee-dao/go-sdk/core"
	"github.com/wetee-dao/libos-entry/util"
)

// 创建一个专门用于为外接用于证明和获取证明的服务
func startEntryServer(fs util.Fs, pk *core.Signer, chainAddr string) error {
	router := chi.NewRouter()
	router.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		// 获取 TEE 证书
		param, err := fs.IssueReport(pk, nil)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
		bt, _ := json.Marshal(param)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(bt)
	})

	teeExecutor := &TeeExecutor{
		chainAddr: chainAddr,
		signer:    pk,
		fs:        fs,
	}

	router.HandleFunc("/tee-call", teeExecutor.HandleTeeCall)

	server := &http.Server{Addr: ":65535", Handler: router}
	fmt.Println("Start entry secret listening http://0.0.0.0:65535 ...")
	err := server.ListenAndServe()
	return err
}
