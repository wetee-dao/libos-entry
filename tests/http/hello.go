package main

import (
	"fmt"
	"net/http"

	inkutil "github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/libos-entry/entry/ego"
	"github.com/wetee-dao/libos-entry/util"
)

func main() {
	// init ego
	err := ego.InitEgo()
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/", indexHandler)
	inkutil.LogWithGreen("SERVE", "http://0.0.0.0:8999")
	err = http.ListenAndServe(":8999", nil)
	fmt.Println(err)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头类型为JSON
	w.Header().Set("Content-Type", "application/json")

	// 创建返回的数据
	response := util.TeeCallBack{
		Args: [][]byte{},
	}

	// 将数据转换为JSON格式
	jsonResponse, _ := response.ToJSON()

	// 写入响应
	w.Write(jsonResponse)
}
