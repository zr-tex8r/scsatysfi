// Copyright (c) 2018 Takayuki YATO (aka. "ZR")
//   GitHub:   https://github.com/zr-tex8r
//   Twitter:  @zr_tex8r
// Distributed under the MIT License.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type argSpec int

const (
	argVoid = iota
	argStr
	argBool
	argBoolOld
	argInt
)

type argInfo struct {
	key  string
	spec argSpec
	proc argOptProc
	doc  string
}

type argProc func(string) error

type argOptProc func(string, string) error

func argSetStr(vp *string) argOptProc {
	return func(arg string, opt string) (err error) {
		*vp = arg
		return
	}
}

func argSetBool(vp *bool) argOptProc {
	return func(arg string, opt string) (err error) {
		*vp = true
		return
	}
}

func argSetInt(vp *int64) argOptProc {
	return func(arg string, opt string) (err error) {
		*vp, err = readInt(arg)
		if err != nil {
			err = fmt.Errorf("option '%s' expects an integer", opt)
		}
		return
	}
}

func argProgramName() string {
	s := os.Args[0]
	return filepath.Base(s[:len(s)-len(filepath.Ext(s))])
}

func argUsageMsg(infos []argInfo) string {
	sb := new(bytes.Buffer)
	sb.WriteString("\n")
	for _, info := range infos {
		sb.WriteString(fmt.Sprintf("  %s %s\n", info.key, info.doc))
	}
	sb.WriteString("  -help  Display this list of options\n")
	sb.WriteString("  --help  Display this list of options\n")
	return sb.String()
}

func argKeyValue(tok string) (key, value string, ok bool) {
	if a := strings.SplitN(tok, "=", 2); len(a) == 2 {
		return a[0], a[1], true
	} else {
		return tok, "", false
	}
}

func argParse(infos []argInfo, anonProc argProc) {
	usage := argUsageMsg(infos)
	for i, lmt := 1, len(os.Args); i < lmt; i++ {
		tok, done := os.Args[i], false
		key, val, valok := argKeyValue(tok)
		if tok == "-help" || tok == "--help" {
			//NB: '-help=VAL' does not count
			fmt.Fprint(os.Stderr, usage)
			os.Exit(0)
		} else if strings.HasPrefix(key, "-") {
			for _, info := range infos {
				if info.key == key {
					switch info.spec {
					case argVoid:
						//NB: value is totally ignored
						if err := info.proc("", key); err != nil {
							argError(usage, "%v", err)
						}
					case argStr, argInt:
						if !valok {
							if i++; i == lmt {
								argError(usage, "'%s' needs an argument", key)
							}
							val = os.Args[i]
						}
						if err := info.proc(val, key); err != nil {
							argValueError(usage, val, err)
						}
					case argBool:
						//NB: value is totally ignored
						if err := info.proc("", key); err != nil {
							argError(usage, "%v", err)
						}
					case argBoolOld: // old behabior of Arg.parse??
						//NB: value is forbidden
						if valok {
							err := fmt.Errorf("option '%s' expects no argument", tok)
							argValueError(usage, val, err)
						}
						if err := info.proc("", key); err != nil {
							argError(usage, "%v", err)
						}
					}
					done = true
					break
				}
			}
		} else {
			if err := anonProc(key); err != nil {
				argError(usage, "%v", err)
			}
			done = true
		}
		if !done {
			argError(usage, "unknown option '%s'", tok)
		}
	}

}

func argValueError(usage, val string, err error) {
	argError(usage, "wrong argument '%s'; %v", val, err)
}

func argError(usage, f string, arg ...interface{}) {
	s := fmt.Sprintf(f, arg...)
	fmt.Fprint(os.Stderr, argProgramName(), ": ", s, ".\n", usage)
	os.Exit(2)
}

//--------readInt

var errBadInt = errors.New("bad integer form")

func readInt(s string) (int64, error) {
	s = strings.Trim(s, " \t\n\f\r")
	// This is an emulation of OCaml's int_of_string,
	// where int has 63 bits.
	neg, sp, base := false, 0, 10
	switch s[0] {
	case '+':
		sp = 1
	case '-':
		neg, sp = true, 1
	}
	if s[sp] == '0' {
		switch s[sp+1] {
		case 'b', 'B':
			sp += 2
			base = 2
		case 'o', 'O':
			sp += 2
			base = 8
		case 'u':
			sp += 2
		case 'x', 'X':
			sp += 2
			base = 16
		}
	}
	if s[sp] == '_' || s[len(s)-1] == '_' {
		return 0, errBadInt
	}
	s = strings.ReplaceAll(s[sp:], "_", "")
	v, err := strconv.ParseInt(s, base, 64)
	if err != nil {
		return 0, err
	}
	if base == 10 && sp < 2 { // no prefix
		if neg {
			v = -v
		}
		if v < -(1<<62) || (1<<62) <= v {
			return 0, errBadInt
		}
	} else {
		if v >= (1 << 62) {
			v = int64(uint64(v) | (1 << 63))
		}
		if neg {
			v = -v
		}
	}
	return v, nil
}
