package main

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/afero"
	"golang.org/x/sys/unix"

	chain "github.com/wetee-dao/ink.go"
	inkutil "github.com/wetee-dao/ink.go/util"
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
	inkutil.LogWithGray("Starting service ", strings.Join(os.Args, " "))
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

// VerifyReport 函数接收一个 *util.TeeParam 类型指针 workerReport
// 并返回一个 *util.TeeReport 类型的指针和一个空的错误
func (l *LibosFs) VerifyReport(workerReport *util.TeeParam) (*util.TeeReport, error) {
	return &util.TeeReport{
		TeeType:       workerReport.TeeType,
		CodeSigner:    []byte{},
		CodeSignature: []byte{},
		CodeProductID: []byte{},
	}, nil
}

// IssueReport 函数为给定的签名者和数据生成一个 TeeParam 报告对象
// 参数：
//
//	pk：签名者对象，包含地址信息
//	data：要包含在报告中的数据
//
// 返回值：
//
//	*util.TeeParam：包含地址、时间、类型、数据和空报告字段的 TeeParam 对象指针
//	error：如果发生错误，返回错误；如果成功，返回 nil
func (l *LibosFs) IssueReport(pk chain.SignerType, data []byte) (*util.TeeParam, error) {
	timestamp := time.Now().Unix()
	return &util.TeeParam{
		Address: pk.Public(),
		Time:    timestamp,
		TeeType: 1,
		Data:    data,
		Report:  []byte{},
	}, nil
}
