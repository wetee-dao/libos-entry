package model

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/edgelesssys/ego/attestation"
	"github.com/edgelesssys/ego/attestation/tcbstatus"
	"github.com/edgelesssys/ego/enclave"
	"github.com/vedhavyas/go-subkey/v2/ed25519"
	chain "github.com/wetee-dao/ink.go"
	"golang.org/x/crypto/blake2b"
)

// sgx issue report
func SgxIssue(pk *chain.Signer, call *TeeCall) error {
	timestamp := time.Now().Unix()
	var buf bytes.Buffer
	buf.Write(Int64ToBytes(timestamp))
	buf.Write(pk.PublicKey)
	if call.Tx != nil {
		txbuf := make([]byte, call.Tx.Size())
		call.Tx.MarshalTo(txbuf)
		buf.Write(txbuf)
	}

	sig, err := pk.Sign(buf.Bytes())
	if err != nil {
		return err
	}

	report, err := enclave.GetRemoteReport(sig)
	if err != nil {
		return err
	}

	// add report to call
	call.Time = timestamp
	call.TeeType = 0
	call.Report = report
	call.Caller = pk.PublicKey

	return nil
}

// sgx verify
func SgxVerify(reportData *TeeCall) (*TeeVerifyResult, error) {
	payload := reportData.Tx
	msgBytes := make([]byte, payload.Size())
	payload.MarshalTo(msgBytes)
	reportBytes, timestamp := reportData.Report, reportData.Time

	report, err := enclave.VerifyRemoteReport(reportBytes)
	if err == attestation.ErrTCBLevelInvalid {
		fmt.Printf("Warning: TCB level is invalid: %v\n%v\n", report.TCBStatus, tcbstatus.Explain(report.TCBStatus))
	} else if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.Write(Int64ToBytes(timestamp))
	buf.Write(reportData.Caller)
	if len(msgBytes) > 0 {
		buf.Write(msgBytes)
	}

	sig := report.Data
	if !SignVerify(reportData.Caller, buf.Bytes(), sig) {
		return nil, errors.New("invalid report sign")
	}

	// if report.Debug {
	// 	return nil, errors.New("debug mode is not allowed")
	// }

	return &TeeVerifyResult{
		TeeType:       reportData.TeeType,
		CodeSigner:    report.SignerID,
		CodeSignature: report.UniqueID,
		CodeProductId: report.ProductID,
	}, nil
}

// client sgx verify
func ClientSgxVerify(reportData *TeeCall) (*TeeVerifyResult, error) {
	payload := reportData.Tx
	msgBytes := make([]byte, payload.Size())
	payload.MarshalTo(msgBytes)
	var reportBytes, timestamp = reportData.Report, reportData.Time

	// decode address
	signer := reportData.Caller

	// call sgx-verify
	reportBt := base64.StdEncoding.EncodeToString(reportBytes)
	cmd := exec.Command("/usr/local/bin/sgx-verify", reportBt)
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.New("call sgx-verify errors" + err.Error())
	}
	outs := strings.Split(string(output), "⊂⊂⊂⊂")
	datas := strings.Split(outs[1], "∐∐∐∐")
	report := &attestation.Report{}
	err = json.Unmarshal([]byte(datas[1]), report)
	if err != nil {
		return nil, errors.New("call sgx-verify unmarshal errors" + err.Error())
	}
	// end call sgx-verify

	var buf bytes.Buffer
	buf.Write(Int64ToBytes(timestamp))
	buf.Write(signer)
	if len(msgBytes) > 0 {
		buf.Write(msgBytes)
	}

	sig := report.Data
	if !SignVerify(reportData.Caller, buf.Bytes(), sig) {
		return nil, errors.New("invalid report sign")
	}

	// if report.Debug {
	// 	return nil, errors.New("debug mode is not allowed")
	// }

	return &TeeVerifyResult{
		TeeType:       reportData.TeeType,
		CodeSigner:    report.SignerID,
		CodeSignature: report.UniqueID,
		CodeProductId: report.ProductID,
	}, nil
}

func SignVerify(pubkeyBt []byte, msg []byte, signature []byte) bool {
	pubkey, err := ed25519.Scheme{}.FromPublicKey(pubkeyBt)
	if err != nil {
		return false
	}

	if len(msg) > 256 {
		h := blake2b.Sum256(msg)
		msg = h[:]
	}

	return pubkey.Verify(msg, signature)
}
