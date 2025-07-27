// Copyright (c) Edgeless Systems GmbH.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package util

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"time"

	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/libos-entry/model"
)

// GramineQuoteIssuer issues quotes.
type GramineQuoteIssuer struct{}

// Issue issues a quote for remote attestation for a given message (usually a certificate).
func (i GramineQuoteIssuer) Issue(pk chain.Signer, call *model.TeeCall) error {
	// hash := sha256.Sum256(cert)
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

	f, err := os.OpenFile("/dev/attestation/user_report_data", os.O_WRONLY, 0)
	if err != nil {
		return err
	}

	_, err = f.Write(sig)
	f.Close()
	if err != nil {
		return err
	}

	f, err = os.Open("/dev/attestation/quote")
	if err != nil {
		return err
	}

	quote := make([]byte, 8192)
	quoteSize, err := f.Read(quote)
	f.Close()
	if err != nil {
		return err
	}

	if !(0 < quoteSize && quoteSize < len(quote)) {
		return errors.New("invalid quote size")
	}

	report := prependOEHeaderToRawQuote(quote[:quoteSize])

	// add report to call
	call.Time = timestamp
	call.TeeType = 0
	call.Report = report
	call.Caller = pk.PublicKey

	return nil
}

func prependOEHeaderToRawQuote(rawQuote []byte) []byte {
	quoteHeader := make([]byte, 16)
	binary.LittleEndian.PutUint32(quoteHeader, 1)     // version
	binary.LittleEndian.PutUint32(quoteHeader[4:], 2) // OE_REPORT_TYPE_SGX_REMOTE
	binary.LittleEndian.PutUint64(quoteHeader[8:], uint64(len(rawQuote)))
	return append(quoteHeader, rawQuote...)
}
