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
	lno, csrlen, ssrc := 0, 0, bufio.NewScanner(rsrc)
	vt, curvt := scNix, scNix
	for ssrc.Scan() {
		lno += 1
		if vt, csrlen, err = scParseLine(lno, csrlen, ssrc.Text()); err != nil {
			return
		}
		if vt == scEssential {
			curvt = scEssential
		}
	}
	if err = ssrc.Err(); err != nil {
		return
	}
	if csrlen > 0 {
		err = sceBadCommentError(lno + 1)
		return
	}
	value = scValue{curvt}
	return
}

func scParseLine(lno, csrlen int, line string) (vt scVType, rcsrlen int, err error) {
	vt = scNix
	srlen := 0
	onSRTerm := func() {
		if csrlen == 0 {
			csrlen = srlen
		} else if csrlen == srlen {
			csrlen = 0
		} // else no-op
		srlen = 0
	}
	for i, r := range line {
		if r == '@' || r == '\U0001F363' { // SUSHI
			srlen += 1
			continue
		} else if srlen > 0 { // sushi-run terminates
			onSRTerm()
		}
		if csrlen > 0 { // in block comment
			continue
		}
		switch r {
		case ' ', '\t':
			// nop
		case '8', '\u2603', '\u26C4', '\u26C7': // SNOWMAN
			vt = scEssential
		case '2', '\U0001F986': // DUCK
			return
		default:
			err = sceBadCharError(lno, i, i+1, r)
			return
		}
	}
	if srlen > 0 { // sushi-run at line end
		onSRTerm()
	}
	rcsrlen = csrlen
	return
}
