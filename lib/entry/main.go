package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/edgelesssys/ego/attestation"
	"github.com/edgelesssys/ego/attestation/tcbstatus"
	"github.com/edgelesssys/ego/eclient"
	"github.com/spf13/afero"
	"golang.org/x/sys/unix"

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
	util.LogWithRed("OS hostfs: ", hostfs)

	var service string

	switch libOS {
	case "Gramine":
		log.Println("Geted libOS: Gramine")

		service, err = libos.InitGramineEntry("", hostfs)
		if err != nil {
			util.ExitWithMsg("Activating Gramine entry failed: %s", err)
		}

		// case "Occlum":
		// 	log.Println("Geted libOS: Occlum")
		// 	service, err =  libos.InitOcclumEntry(hostfs)
		// 	if err != nil {
		// 		exit("activating Occlum entry failed: %s", err)
		// 	}
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

func (l *LibosFs) VerifyReport(reportBytes, certBytes, signer []byte) error {
	report, err := eclient.VerifyRemoteReport(reportBytes)
	if err == attestation.ErrTCBLevelInvalid {
		fmt.Printf("Warning: TCB level is invalid: %v\n%v\n", report.TCBStatus, tcbstatus.Explain(report.TCBStatus))
		fmt.Println("We'll ignore this issue in this sample. For an app that should run in production, you must decide which of the different TCBStatus values are acceptable for you to continue.")
	} else if err != nil {
		return err
	}

	hash := sha256.Sum256(certBytes)
	if !bytes.Equal(report.Data[:len(hash)], hash[:]) {
		return errors.New("report data does not match the certificate's hash")
	}

	// You can either verify the UniqueID or the tuple (SignerID, ProductID, SecurityVersion, Debug).

	if report.SecurityVersion < 2 {
		return errors.New("invalid security version")
	}
	if binary.LittleEndian.Uint16(report.ProductID) != 1234 {
		return errors.New("invalid product")
	}
	if !bytes.Equal(report.SignerID, signer) {
		return errors.New("invalid signer")
	}

	// For production, you must also verify that report.Debug == false

	return nil
}

func (l *LibosFs) IssueReport(data []byte) ([]byte, error) {
	switch l.LibOsType {
	case "Gramine":
		gramine := util.GramineQuoteIssuer{}
		return gramine.Issue(data)
	default:
		return nil, errors.New("Invalid LibOsType")
	}
}

func (l *LibosFs) SetPassword(password string) {

}
