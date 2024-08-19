package libos

import (
	"encoding/hex"
	"errors"

	"github.com/vedhavyas/go-subkey/v2"
	chain "github.com/wetee-dao/go-sdk"
	"github.com/wetee-dao/go-sdk/module"
	"github.com/wetee-dao/libos-entry/util"
)

// VerifyWorker 函数验证工人报告并返回签名者或错误
func VerifyWorker(reportData *util.TeeParam, fs util.Fs, client *chain.ChainClient) ([]byte, error) {
	// 解码地址
	_, signer, err := subkey.SS58Decode(reportData.Address)
	if err != nil {
		return nil, errors.New("SS58 decode: " + err.Error())
	}

	report, err := VerifyReport(reportData, fs)
	if err != nil {
		return nil, errors.New("verify cluster report: " + err.Error())
	}

	// 校验 worker 代码版本
	codeHash, codeSigner, err := module.GetWorkerCode(client)
	if err != nil {
		return nil, errors.New("GetWorkerCode error:" + err.Error())
	}
	if len(codeHash) > 0 || len(codeSigner) > 0 {
		if hex.EncodeToString(codeHash) != hex.EncodeToString(report.UniqueID) {
			return nil, errors.New("worker code hash error")
		}

		if hex.EncodeToString(codeSigner) != hex.EncodeToString(report.SignerID) {
			return nil, errors.New("worker signer error")
		}
	}

	return signer, nil
}
