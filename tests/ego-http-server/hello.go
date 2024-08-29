package main

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
	"github.com/wetee-dao/libos-entry/entry/ego"
	"github.com/wetee-dao/libos-entry/util"
)

func main() {
	err := ego.InitEgo()
	if err != nil {
		fmt.Println(err)
		return
	}
	http.HandleFunc("/", indexHandler)
	err = http.ListenAndServe(":8999", nil)
	fmt.Println(err)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头类型为JSON
	w.Header().Set("Content-Type", "application/json")

	// 创建返回的数据
	response := util.TeeCallBack{
		Args: []gtypes.InkArg{
			{
				IsU128:       true,
				AsU128Field0: types.NewU128(*big.NewInt(500)),
			},
			{
				IsBool:       true,
				AsBoolField0: true,
			},
			{
				IsTString:       true,
				AsTStringField0: []byte("Hello, World!"),
			},
		},
	}

	// 将数据转换为JSON格式
	jsonResponse, _ := response.ToJSON()

	// 写入响应
	w.Write(jsonResponse)
}
