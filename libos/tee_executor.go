package libos

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/go-resty/resty/v2"
	"github.com/vedhavyas/go-subkey/v2"
	"github.com/wetee-dao/go-sdk/core"
	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
	"github.com/wetee-dao/libos-entry/libos/chain"
	"github.com/wetee-dao/libos-entry/util"
)

type TeeExecutor struct {
	chainAddr string
	signer    *core.Signer
	fs        util.Fs
}

func (e *TeeExecutor) HandleTeeCall(w http.ResponseWriter, r *http.Request) {
	/// 解析请求
	t := util.TeeTrigger{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		r.Body.Close()
		fmt.Println("io.ReadAll", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	err = json.Unmarshal(body, &t)
	if err != nil {
		fmt.Println("json.Unmarshal(", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	workerReport := t.Tee
	msg := fmt.Sprint(t.ClusterId, t.Callids)

	// decode address
	_, signer, err := subkey.SS58Decode(t.Tee.Address)
	if err != nil {
		fmt.Println("subkey.SS58Decode ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取 worker report
	_, err = e.fs.VerifyReport(workerReport.Report, []byte(msg), signer, workerReport.Time)
	if err != nil {
		fmt.Println("VerifyReport", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = e.runCallAndSubmit(&t)
	if err != nil {
		fmt.Println("RunCallAndSubmit", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

func (e *TeeExecutor) runCallAndSubmit(t *util.TeeTrigger) error {
	// 初始化去快链链接
	c, err := chain.InitChain(e.chainAddr, e.signer)
	if err != nil {
		return errors.New("chain.InitChain: " + err.Error())
	}
	defer c.Close()

	callIds := make([]types.U128, 0, len(t.Callids))
	for _, callid := range t.Callids {
		n := new(big.Int)
		n.SetString(callid, 10)
		callIds = append(callIds, types.NewU128(*n))
	}

	// 获取 tee calls
	calls, _, err := c.ListTeeCalls(t.ClusterId, callIds)
	if err != nil {
		return err
	}

	if len(calls) == 0 {
		return errors.New("no tee calls")
	}

	// 获取 meta api
	meta, err := c.GetMetaApi(calls[0].WorkId)
	if err != nil {
		return err
	}

	// 运行 tee calls
	resps := make([]chain.TeeCallBack, 0, len(calls))
	for _, call := range calls {
		resp, err := callTeeApp(call, &meta)
		resps = append(resps, chain.TeeCallBack{
			Err:  err,
			Resp: resp,
		})
	}

	// 提交 proof
	err = c.TeeCallback(t.ClusterId, callIds, resps)
	if err != nil {
		return err
	}
	fmt.Println("runCallAndSubmit", "success")

	return nil
}

// CallTeeApp 调用 tee app api
func callTeeApp(call *gtypes.TEECall, meta *gtypes.ApiMeta) ([]byte, error) {
	client := resty.New()
	req := client.R().SetBody(call.Args)

	// 构造请求参数
	api := meta.Apis[call.Method]
	url := "http://0.0.0.0:" + fmt.Sprint(meta.Port) + string(api.Url)

	// 0: get, 1: post, 2: put, 3: delete
	switch api.Method {
	case 0:
		resp, err := req.Get(url)
		if err != nil {
			return nil, err
		}
		return resp.Body(), nil
	case 1:
		resp, err := req.Post(url)
		if err != nil {
			return nil, err
		}
		return resp.Body(), nil
	case 2:
		resp, err := req.Put(url)
		if err != nil {
			return nil, err
		}
		return resp.Body(), nil
	case 3:
		resp, err := req.Delete(url)
		if err != nil {
			return nil, err
		}
		return resp.Body(), nil
	default:
		return nil, nil
	}
}
