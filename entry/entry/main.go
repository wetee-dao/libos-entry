package main

import (
	"errors"
	"log"
	"os"
	"strings"

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

	switch libOS {
	case "Gramine":
		log.Println("Geted libOS: Gramine")

		service, err = libos.InitGramineEntry(hostfs, false)
		if err != nil {
			util.ExitWithMsg("Activating Gramine entry failed: %s", err)
		}
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
	gramine   *util.GramineQuoteIssuer
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
	// report, err := eclient.VerifyRemoteReport(reportBytes)
	// if err == attestation.ErrTCBLevelInvalid {
	// 	fmt.Printf("Warning: TCB level is invalid: %v\n%v\n", report.TCBStatus, tcbstatus.Explain(report.TCBStatus))
	// 	fmt.Println("We'll ignore this issue in this sample. For an app that should run in production, you must decide which of the different TCBStatus values are acceptable for you to continue.")
	// } else if err != nil {
	// 	return err
	// }

	return &util.TeeReport{
		TeeType:       workerReport.TeeType,
		CodeSigner:    []byte{},
		CodeSignature: []byte{},
		CodeProductID: []byte{},
	}, nil
}

func (l *LibosFs) IssueReport(pk chain.SignerType, data []byte) (*util.TeeParam, error) {
	switch l.LibOsType {
	case "Gramine":
		if l.gramine == nil {
			l.gramine = &util.GramineQuoteIssuer{}
		}
		report, t, err := l.gramine.Issue(pk, data)
		if err != nil {
			return nil, err
		}
		return &util.TeeParam{
			Time: t,
			// Address: pk.SS58Address(42),
			Report: report,
			Data:   data,
		}, nil
	default:
		return nil, errors.New("invalid libos")
	}
}
