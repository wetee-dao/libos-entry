package model

import (
	"bytes"
	"errors"
	"time"

	"github.com/cometbft/cometbft/abci/types"
	"github.com/google/go-sev-guest/client"
	"github.com/google/go-sev-guest/proto/sevsnp"
	"github.com/google/go-sev-guest/verify"
	"github.com/vedhavyas/go-subkey/v2/ed25519"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/libos-entry/model/protoio"
)

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
	defer device.Close()

	sig64 := *(*[64]byte)(sig[:64])
	attestationReport, err := client.GetExtendedReport(device, sig64)
	if err != nil {
		return errors.New("client.GetExtendedReport:" + err.Error())
	}

	reportBuf := new(bytes.Buffer)
	err = types.WriteMessage(attestationReport, reportBuf)
	if err != nil {
		return err
	}

	// add report to call
	call.Time = timestamp
	call.TeeType = 1
	call.Report = reportBuf.Bytes()
	call.Caller = pk.PublicKey

	return nil
}

func SnpVerify(reportData *TeeCall) (*TeeVerifyResult, error) {
	payload := reportData.Tx
	msgBytes := make([]byte, payload.Size())
	payload.MarshalTo(msgBytes)
	var reportBytes, timestamp = reportData.Report, reportData.Time

	signer := reportData.Caller

	attestation := new(sevsnp.Attestation)
	err := protoio.ReadMessage(bytes.NewBuffer(reportBytes), attestation)
	if err != nil {
		return nil, err
	}

	// 验证报告
	options := verify.DefaultOptions()
	if err = verify.SnpAttestation(attestation, options); err != nil {
		return nil, errors.New("verify.SnpAttestation error:" + err.Error())
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

	sig := attestation.Report.ReportData
	if !pubkey.Verify(buf.Bytes(), sig) {
		return nil, errors.New("invalid sgx report")
	}

	return &TeeVerifyResult{
		TeeType:       reportData.TeeType,
		CodeSigner:    []byte{},
		CodeSignature: []byte{},
		CodeProductId: []byte{},
	}, nil
}
