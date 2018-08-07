// Copyright (c) 2018 Takayuki YATO (aka. "ZR")
//   GitHub:   https://github.com/zr-tex8r
//   Twitter:  @zr_tex8r
// Distributed under the MIT License.

package main

import (
	"bufio"
	"io"
)

//--------scVType

type scVType int

const (
	scNix = scVType(iota)
	scEssential
)

var scVTypeName = []string{"nix", "essential"}

func (t scVType) String() string {
	return scVTypeName[t]
}

//--------scValue

type scValue struct {
	vtype scVType
	// NB: no other data is needed
}

//--------scParse

func scParseReader(rsrc io.Reader) (value scValue, err error) {
	lno, ssrc := 0, bufio.NewScanner(rsrc)
	vt, curvt := scNix, scNix
	for ssrc.Scan() {
		lno += 1
		if vt, err = scParseLine(lno, ssrc.Text()); err != nil {
			return
		}
		if vt == scEssential {
			curvt = scEssential
		}
	}
	if err = ssrc.Err(); err != nil {
		return
	}
	value = scValue{curvt}
	return
}

func scParseLine(lno int, line string) (vt scVType, err error) {
	vt = scNix
	for i, r := range line {
		switch r {
		case ' ', '\t':
			// nop
		case '8', '\u2603', '\u26C4', '\u26C7':
			vt = scEssential
		case '2', '\U0001F986':
			return
		default:
			err = sceBadCharError(lno, i, i+1, r)
			return
		}
	}
	return
}
