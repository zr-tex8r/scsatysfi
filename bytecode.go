// Copyright (c) 2018 Takayuki YATO (aka. "ZR")
//   GitHub:   https://github.com/zr-tex8r
//   Twitter:  @zr_tex8r
// Distributed under the MIT License.

package main

import (
	"bytes"
	"errors"
)

// This is the only valid code string.
var codeDoEssential = []byte{0x38}

type scByteCode struct {
	code []byte
}

func byteCompile(value scValue) scByteCode {
	// Since the value is type-checked, it never fails.
	return scByteCode{codeDoEssential}
}

func byteExec(bcode scByteCode) {
	if bytes.Equal(bcode.code, codeDoEssential) {
		writeOutput(outFile)
	} else {
		scePanic(errors.New("bad bytecode"))
	}
}
