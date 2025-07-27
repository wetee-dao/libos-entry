package libos

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	chain "github.com/wetee-dao/ink.go"
	inkutil "github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/libos-entry/util"
)

// 创建一个专门用于为外接用于证明和获取证明的服务
func startTEEServer(fs util.Fs, pk chain.SignerType, chainAddr string) error {
	router := chi.NewRouter()
	router.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		// 获取 TEE 证书
		// param, err := fs.IssueReport(pk, nil)
		// if err != nil {
		// 	w.WriteHeader(500)
		// 	w.Write([]byte(err.Error()))
		// }
		// bt, _ := json.Marshal(param)
		bt := []byte("123")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(bt)
	})

	if chainAddr != "" {
		// teeExecutor := &TeeExecutor{
		// 	chainAddr: chainAddr,
		// 	signer:    pk,
		// 	fs:        fs,
		// }
		// router.HandleFunc("/tee-call", teeExecutor.HandleTeeCall)
	}

	server := &http.Server{Addr: ":65535", Handler: router}
	inkutil.LogWithGreen("TEE server", "http://0.0.0.0:65535")
	err := server.ListenAndServe()
	return err
}
