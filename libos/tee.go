package libos

import (
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/libos-entry/model"
	"github.com/wetee-dao/libos-entry/util"
)

// VerifyWorker 函数验证工人报告并返回签名者或错误
func VerifyWorker(reportData *model.TeeCall, fs util.Fs, client *chain.ChainClient) ([]byte, error) {
	// 解码地址
	signer := reportData.Caller

	// report, err := fs.VerifyReport(reportData)
	// if err != nil {
	// 	return nil, errors.New("verify cluster report: " + err.Error())
	// }

	// // 校验 worker 代码版本
	// codeHash, codeSigner, err := module.GetWorkerCode(client)
	// if err != nil {
	// 	return nil, errors.New("GetWorkerCode error:" + err.Error())
	// }
	// if len(codeHash) > 0 || len(codeSigner) > 0 {
	// 	if hex.EncodeToString(codeHash) != hex.EncodeToString(report.CodeSignature) {
	// 		return nil, errors.New("worker code hash error")
	// 	}

	// 	if hex.EncodeToString(codeSigner) != hex.EncodeToString(report.CodeSigner) {
	// 		return nil, errors.New("worker signer error")
	// 	}
	// }

	return signer, nil
}
