// Protocol Buffers for Go with Gadgets
//
// Copyright (c) 2013, The GoGo Authors. All rights reserved.
// http://github.com/gogo/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package protoio

import (
	"encoding/binary"
	"io"

	"github.com/cosmos/gogoproto/proto"
	protov2 "google.golang.org/protobuf/proto"
)

// WriteMessage writes a varint length-delimited protobuf message with deterministic serialization.
// This ensures consistent hash values across multiple serializations.
func WriteMessage(msg proto.Message, w io.Writer) error {
	data, err := MarshalDeterministic(msg)
	if err != nil {
		return err
	}

	var buf [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(buf[:], uint64(len(data)))
	if _, err := w.Write(buf[:n]); err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// MarshalDeterministic serializes a protobuf message with deterministic output.
// This ensures consistent hash values across multiple serializations.
func MarshalDeterministic(msg proto.Message) ([]byte, error) {
	// Try protov2 Message first (newer protobuf)
	if m, ok := any(msg).(protov2.Message); ok {
		return protov2.MarshalOptions{Deterministic: true}.Marshal(m)
	}

	// Fall back to gogoproto Marshal
	// Note: gogoproto's Marshal is deterministic for messages without map fields
	return proto.Marshal(msg)
}

// NewDelimitedWriter returns a new delimited writer that writes varint-delimited
// protobuf messages with deterministic serialization.
func NewDelimitedWriter(w io.Writer) WriteCloser {
	return &varintWriter{w: w}
}

type varintWriter struct {
	w io.Writer
}

func (w *varintWriter) WriteMsg(msg proto.Message) (int, error) {
	data, err := MarshalDeterministic(msg)
	if err != nil {
		return 0, err
	}

	var buf [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(buf[:], uint64(len(data)))
	written, err := w.w.Write(buf[:n])
	if err != nil {
		return written, err
	}

	nw, err := w.w.Write(data)
	written += nw
	return written, err
}

func (w *varintWriter) Close() error {
	if c, ok := w.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
