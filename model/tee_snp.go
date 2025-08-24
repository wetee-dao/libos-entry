package model

import (
	"bytes"
	"fmt"
	"time"

	"github.com/google/go-sev-guest/client"
	"github.com/google/go-sev-guest/proto/sevsnp"
	"github.com/google/go-sev-guest/verify"
	"github.com/pkg/errors"
	chain "github.com/wetee-dao/ink.go"
	"google.golang.org/protobuf/proto"
)

// snp issue
func SnpIssue(pk *chain.Signer, call *TeeCall) error {
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

	device, err := client.OpenDevice()
	if err != nil {
		return err
	}

	sig64 := *(*[64]byte)(sig[:64])
	attestationReport, err := client.GetExtendedReport(device, sig64)
	if err != nil {
		device.Close()
		return errors.New("client.GetExtendedReport:" + err.Error())
	}
	device.Close()

	data, err := proto.Marshal(attestationReport)
	if err != nil {
		return err
	}

	// add report to call
	call.Time = timestamp
	call.TeeType = 1
	call.Report = data
	call.Caller = pk.PublicKey

	return nil
}

// snp verify
func SnpVerify(callData *TeeCall) (result *TeeVerifyResult, err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			result = nil
			err = errors.New("SnpVerify recover error: " + fmt.Sprint(rerr))
		}
	}()

	reportBytes, timestamp, signer := callData.Report, callData.Time, callData.Caller

	// 解析报告
	attestation := new(sevsnp.Attestation)
	err = proto.Unmarshal(reportBytes, attestation)
	if err != nil {
		return nil, err
	}

	// 验证报告
	options := verify.DefaultOptions()
	if err = verify.SnpAttestation(attestation, options); err != nil {
		return nil, errors.Wrap(err, "verify.SnpAttestation")
	}

	// 构建签名数据
	var buf bytes.Buffer
	buf.Write(Int64ToBytes(timestamp))
	buf.Write(signer)
	payload := callData.Tx
	msgBytes := make([]byte, payload.Size())
	payload.MarshalTo(msgBytes)
	if len(msgBytes) > 0 {
		buf.Write(msgBytes)
	}

	// 验证签名
	sig := attestation.Report.ReportData
	if !SignVerify(callData.Caller, buf.Bytes(), sig) {
		return nil, errors.New("invalid report sign")
	}

	return &TeeVerifyResult{
		TeeType:       callData.TeeType,
		CodeSigner:    []byte{},
		CodeSignature: []byte{},
		CodeProductId: []byte{},
	}, nil
}
