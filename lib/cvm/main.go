package main

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/afero"
	"golang.org/x/sys/unix"

	"github.com/wetee-dao/go-sdk/core"
	"github.com/wetee-dao/libos-entry/libos"
	"github.com/wetee-dao/libos-entry/util"
)

func main() {
	// Get libOS based on uname
	// 获取 libOS 类型
	libOS, err := util.GetLibOS()
	if err != nil {
		util.ExitWithMsg("Failed to get libOS: %s", err)
	}

	// Use filesystem from libOS
	// 获取libOS的文件系统
	hostfs := &LibosFs{LibOsType: libOS}

	var service string
	service, err = libos.InitGramineEntry(hostfs, true)
	if err != nil {
		util.ExitWithMsg("Activating entry failed: %s", err)
	}

	// Start service
	// 开启服务
	util.LogWithRed("Starting service ", strings.Join(os.Args, " "))
	if err := unix.Exec(service, os.Args, os.Environ()); err != nil {
		util.ExitWithMsg("Starting service error", err.Error())
	}
}

type LibosFs struct {
	afero.OsFs
	LibOsType string
}

// Read implements util.Fs.
func (f *LibosFs) ReadFile(filename string) ([]byte, error) {
	return afero.ReadFile(f, filename)
}

// Write implements util.Fs.
func (f *LibosFs) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return afero.WriteFile(f, filename, data, perm)
}

func (l *LibosFs) VerifyReport(workerReport *util.TeeParam) (*util.TeeReport, error) {
	return &util.TeeReport{
		TeeType:       workerReport.TeeType,
		CodeSigner:    []byte{},
		CodeSignature: []byte{},
		CodeProductID: []byte{},
	}, nil
}

func (l *LibosFs) IssueReport(pk *core.Signer, data []byte) (*util.TeeParam, error) {
	timestamp := time.Now().Unix()
	return &util.TeeParam{
		Address: pk.Address,
		Time:    timestamp,
		TeeType: 1,
		Data:    data,
		Report:  []byte{},
	}, nil
}
