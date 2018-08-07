// Copyright (c) 2018 Takayuki YATO (aka. "ZR")
//   GitHub:   https://github.com/zr-tex8r
//   Twitter:  @zr_tex8r
// Distributed under the MIT License.

package main

import (
	"errors"
	"fmt"
	"github.com/zr-tex8r/scpdf"
	"github.com/zr-tex8r/xcolor"
	"image/color"
	"os"
)

const (
	progName = "scSATySFi"
	version  = "0.8.8-pre02"
)

const dfltMuffler = "cmyk:red,1"

var (
	inFile         string
	outFile        string
	fullPath       bool
	mufflerVal     string
	debugShowBbox  bool
	debugShowSpace bool
	typeCheckOnly  bool
	byteComp       bool
	mufflerColor   color.Color
)

func showVersion(string) error {
	fmt.Printf("  %s version %s\n", progName, version)
	os.Exit(0)
	return nil
}

var argSpecList = []argInfo{
	argInfo{"-o", argStr, argSetStr(&outFile), " Specify output file"},
	argInfo{"--output", argStr, argSetStr(&outFile), " Specify output file"},
	argInfo{"-v", argVoid, showVersion, " Print version"},
	argInfo{"--version", argVoid, showVersion, " Print version"},
	argInfo{"--full-path", argBool, argSetBool(&fullPath), " Display paths in full-path style"},
	argInfo{"--debug-show-bbox", argBool, argSetBool(&debugShowBbox), " Output bounding boxes for glyphs"},
	argInfo{"--debug-show-space", argBool, argSetBool(&debugShowSpace), " Output boxes for spaces"},
	argInfo{"-t", argBool, argSetBool(&typeCheckOnly), " Stops after type checking"},
	argInfo{"--type-check-only", argBool, argSetBool(&typeCheckOnly), " Stops after type checking"},
	argInfo{"-b", argBool, argSetBool(&byteComp), " Use bytecode compiler"},
	argInfo{"--bytecomp", argBool, argSetBool(&byteComp), " Use bytecode compiler"},
	argInfo{"--muffler", argStr, argSetStr(&mufflerVal), " Specify muffler color"},
}

func readArg() {
	argParse(argSpecList, func(arg string) error {
		inFile = arg
		return nil
	})
	if inFile == "" {
		scePanic(errors.New("no input file designation."))
	}
	if outFile == "" {
		outFile = changeExt(inFile, ".pdf")
	}
	if mufflerVal == "" {
		mufflerVal = dfltMuffler
	}
	var err error
	if mufflerColor, err = xcolor.GoColor(mufflerVal); err != nil {
		scePanic(err)
	}
}

func main() {
	readArg()

	fmt.Printf(" ---- ---- ---- ----\n")
	fmt.Printf("  target file: '%s'\n", ordPath(outFile))
	aux := changeExt(outFile, ".scsatysfi-aux")
	fmt.Printf("  dump file: '%s' (won't be created)\n", ordPath(aux))

	value := parseFile(inFile)

	readFile(inFile, value)
	if typeCheckOnly {
		return
	}

	fmt.Printf(" ---- ---- ---- ----\n")
	fmt.Printf("  evaluating texts ...\n")
	fmt.Printf("  evaluation done.\n")

	if byteComp {
		byteExec(byteCompile(value))
	} else {
		writePdf(outFile)
	}
}

func writePdf(ppdf string) {
	fmt.Printf(" ---- ---- ---- ----\n")
	fmt.Printf("  writing pages ...\n")

	wpdf, err := os.Create(ppdf)
	sceAssert(err)
	defer wpdf.Close()

	doc := new(scpdf.Doc)
	doc.SetDocInfo(map[string]string{
		"title":   "\u2603",
		"creator": "scSATySFi",
	})
	doc.AddPage(mufflerColor)
	_, err = doc.WriteTo(wpdf)
	sceAssert(err)

	fmt.Printf(" ---- ---- ---- ----\n")
	fmt.Printf("  output written on '%s'.\n", ordPath(ppdf))
}

func readFile(psrc string, value scValue) {
	fmt.Printf(" ---- ---- ---- ----\n")
	fmt.Printf("  reading '%s' ...\n", ordPath(psrc))
	fmt.Printf("  type check passed. (%s)\n", value.vtype)

	if !typeCheckOnly && value.vtype != scEssential {
		scePanic(sceNonDocError(psrc, value.vtype))
	}
}

func parseFile(psrc string) (value scValue) {
	fmt.Printf("  parsing '%s' ...\n", ordPath(psrc))

	rsrc, err := os.Open(psrc)
	sceAssert(err)
	defer rsrc.Close()
	value, err = scParseReader(rsrc)
	sceAssert(err)
	return
}
