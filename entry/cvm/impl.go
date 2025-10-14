package main

import (
	"encoding/json"

	inkutil "github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/libos-entry/libos"
)

type CvmServer struct {
}

func init() {
	TEEServerImpl = CvmServer{}
}

func (CvmServer) start(req *CrossRequest) CrossResponse {
	envs := map[string]string{}
	err := json.Unmarshal(req.env, &envs)
	if err != nil {
		return CrossResponse{code: 1, data: []byte(err.Error())}
	}

	hostfs := &Fs{}
	datas, err := libos.PreLoadFromInitData(hostfs, envs, false)
	if err != nil {
		inkutil.LogWithGray("PreLoadFromInitData", err.Error())
		return CrossResponse{code: 1, data: []byte(err.Error())}
	}

	bt, _ := json.Marshal(datas)
	return CrossResponse{code: 0, images: libos.Images, data: bt}
}

// stop implements TEEServer.
func (c CvmServer) stop(_ *int64) CrossResponse {
	panic("unimplemented")
}
