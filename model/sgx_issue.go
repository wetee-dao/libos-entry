package model

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/edgelesssys/ego/attestation"
	"github.com/edgelesssys/ego/attestation/tcbstatus"
	"github.com/edgelesssys/ego/enclave"
	"github.com/vedhavyas/go-subkey/v2/ed25519"
	chain "github.com/wetee-dao/ink.go"
)

func IssueReport(pk *chain.Signer, call *TeeCall) error {
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

func VerifyReport(reportData *TeeCall) (*TeeVerifyResult, error) {
	// TODO SEV/TDX not support
	if reportData.TeeType != 0 {
		return &TeeVerifyResult{
			CodeSignature: []byte{},
			CodeSigner:    []byte{},
			CodeProductId: []byte{},
		}, nil
	}

	payload := reportData.Tx
	msgBytes := make([]byte, payload.Size())
	payload.MarshalTo(msgBytes)
	var reportBytes, timestamp = reportData.Report, reportData.Time

	// decode address
	signer := reportData.Caller

	report, err := enclave.VerifyRemoteReport(reportBytes)
	if err == attestation.ErrTCBLevelInvalid {
		fmt.Printf("Warning: TCB level is invalid: %v\n%v\n", report.TCBStatus, tcbstatus.Explain(report.TCBStatus))
	} else if err != nil {
		return nil, err
	}

	pubkey, err := ed25519.Scheme{}.FromPublicKey(signer)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.Write(Int64ToBytes(timestamp))
	buf.Write(signer)
	if len(msgBytes) > 0 {
		buf.Write(msgBytes)
	}

	sig := report.Data
	if !pubkey.Verify(buf.Bytes(), sig) {
		return nil, errors.New("invalid sgx report")
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

func Int64ToBytes(time int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(time))
	return b
}

func BytesToInt64(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}
