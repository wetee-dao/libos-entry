package chain

import (
	"fmt"
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/go-resty/resty/v2"
	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
	"github.com/wetee-dao/go-sdk/pallet/weteebridge"
)

// HandlerTeeCall 处理 tee 调用
func (c *Chain) HandlerTeeCall(call *gtypes.TEECall, meta *gtypes.ApiMeta) {
	// 调用 tee app
	// call tee app
	resp, err := CallTeeApp(call, meta)
	if err != nil {
		return
	}

	// 回调结果到区块链
	// callback to chain
	recall := weteebridge.MakeInkCallbackCall(1, call.Id, resp, types.NewU128(*big.NewInt(0)))
	err = c.client.SignAndSubmit(c.signer, recall, false)
	if err != nil {
		fmt.Println("callback to chain error:", err)
	} else {
		fmt.Println("callback to chain success")
	}
}

// CallTeeApp 调用 tee app api
func CallTeeApp(call *gtypes.TEECall, meta *gtypes.ApiMeta) ([]byte, error) {
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
