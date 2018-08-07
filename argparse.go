// Copyright (c) 2018 Takayuki YATO (aka. "ZR")
//   GitHub:   https://github.com/zr-tex8r
//   Twitter:  @zr_tex8r
// Distributed under the MIT License.

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type argSpec int

const (
	argVoid = iota
	argStr
	argBool
	argBoolOld
)

type argInfo struct {
	key  string
	spec argSpec
	proc argProc
	doc  string
}

type argProc func(string) error

func argSetStr(vp *string) argProc {
	return func(arg string) (err error) {
		*vp = arg
		return
	}
}

func argSetBool(vp *bool) argProc {
	return func(arg string) (err error) {
		*vp = true
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
						if err := info.proc(""); err != nil {
							argError(usage, "%v", err)
						}
					case argStr:
						if !valok {
							if i++; i == lmt {
								argError(usage, "'%s' needs an argument", key)
							}
							val = os.Args[i]
						}
						if err := info.proc(val); err != nil {
							argValueError(usage, val, err)
						}
					case argBool:
						//NB: value is totally ignored
						if err := info.proc(""); err != nil {
							argError(usage, "%v", err)
						}
					case argBoolOld: // old behabior of Arg.parse??
						//NB: value is forbidden
						if valok {
							err := fmt.Errorf("option '%s' expects no argument", tok)
							argValueError(usage, val, err)
						}
						if err := info.proc(""); err != nil {
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
