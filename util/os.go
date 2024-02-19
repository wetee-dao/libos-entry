// Copyright (c) Edgeless Systems GmbH.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/sys/unix"
)

func GetLibOS() (string, error) {
	utsname := unix.Utsname{}
	if err := unix.Uname(&utsname); err != nil {
		return "NotTEE", err
	}

	// Clean utsname
	sysname := strings.ReplaceAll(string(utsname.Sysname[:]), "\x00", "")
	release := strings.ReplaceAll(string(utsname.Release[:]), "\x00", "")
	version := strings.ReplaceAll(string(utsname.Version[:]), "\x00", "")
	machine := strings.ReplaceAll(string(utsname.Machine[:]), "\x00", "")

	// Occlum
	// Taken from: https://github.com/occlum/occlum/blob/master/src/libos/src/misc/uname.rs
	if sysname == "Occlum" {
		return "Occlum", nil
	}

	// Gramine
	// This looks like a general Linux kernel name, making it harder to get... But it's unlikely someone is running SGX code on Linux 3.10.0.
	// Taken from: https://github.com/gramineproject/gramine/blob/c83ec08f10cdbb3a258d18b118dd95602a55abc9/libos/src/sys/libos_uname.c
	if sysname == "Linux" && release == "3.10.0" && version == "1" && machine == "x86_64" {
		return "Gramine", nil
	}

	return "Gramine", nil
	// return "NotTEE", errors.New("cannot get libOS")
}

func ExitWithMsg(format string, args ...interface{}) {
	// Print error message in red and append newline
	// then exit with error code 1
	msg := fmt.Sprintf("Error: %s\n", format)
	_, _ = color.New(color.FgRed).Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}
