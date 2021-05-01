// Copyright (c) 2018-2020 Takayuki YATO (aka. "ZR")
//   GitHub:   https://github.com/zr-tex8r
//   Twitter:  @zr_tex8r
// Distributed under the MIT License.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/zr-tex8r/scpdf"
	"github.com/zr-tex8r/xcolor"
)

const (
	progName = "scSATySFi"
	version  = "0.8.58"
)

const dfltMuffler = "cmyk:red,1"

var (
	inFile              string
	outFile             string
	fullPath            bool
	mufflerVal          string
	debugShowBbox       bool
	debugShowSpace      bool
	debugShowBlockBbox  bool
	debugShowBlockSpace bool
	debugShowOverfull   bool
	typeCheckOnly       bool
	byteComp            bool
	mufflerColor        color.Color
	textModeVal         string
	textMode            string
	markdownVal         string
	isMarkdown          bool
	evalVal             string
	isShowFont          bool
	config              string
	noDefaultConfig     bool
	pageNumberLimit     string
)

func showVersion(string) error {
	fmt.Printf("  %s version %s\n", progName, version)
	os.Exit(0)
	return nil
}

var argSpecList = []argInfo{
	argInfo{"-o", argStr, argSetStr(&outFile), " Specify output file"},
	argInfo{"--output", argStr, argSetStr(&outFile), " Specify output file"},
	argInfo{"-v", argVoid, showVersion, " Prints version"},
	argInfo{"--version", argVoid, showVersion, " Prints version"},
	argInfo{"--full-path", argBool, argSetBool(&fullPath), " Displays paths in full-path style"},
	// But there is no such inessential thing as glyph
	argInfo{"--debug-show-bbox", argBool, argSetBool(&debugShowBbox), " Outputs bounding boxes for glyphs"},
	// But there is no such inessential thing as space
	argInfo{"--debug-show-space", argBool, argSetBool(&debugShowSpace), " Outputs boxes for spaces"},
	// Again...
	argInfo{"--debug-show-block-bbox", argBool, argSetBool(&debugShowBlockBbox), " Outputs bounding boxes for blocks"},
	argInfo{"--debug-show-block-space", argBool, argSetBool(&debugShowBlockSpace), " Outputs visualized block spaces"},
	argInfo{"--debug-show-overfull", argBool, argSetBool(&debugShowOverfull), " Outputs visualized overfull or underfull lines"},
	argInfo{"-t", argBool, argSetBool(&typeCheckOnly), " Stops after type checking"},
	argInfo{"--type-check-only", argBool, argSetBool(&typeCheckOnly), " Stops after type checking"},
	argInfo{"-b", argBool, argSetBool(&byteComp), " Use bytecode compiler"},
	argInfo{"--bytecomp", argBool, argSetBool(&byteComp), " Use bytecode compiler"},
	argInfo{"--text-mode", argStr, argSetStr(&textModeVal), " Set text mode"},
	argInfo{"--markdown", argStr, argSetStr(&markdownVal), " Pass Markdown source as input"},
	// But there is no such inessential thing as font
	argInfo{"--show-fonts", argBool, argSetBool(&isShowFont), " Displays all the available fonts"},
	// But there is no such inessential thing as config
	argInfo{"-C", argStr, argSetStr(&config), " Add colon-separated paths to configuration search path"},
	argInfo{"--config", argStr, argSetStr(&config), " Add colon-separated paths to configuration search path"},
	// Again...
	argInfo{"--no-default-config", argBool, argSetBool(&noDefaultConfig), " Does not use default configuration search path"},
	// How does it work?
	argInfo{"--page-number-limit", argStr, argSetStr(&pageNumberLimit), " Set the page number limit (default: 10000)"},
	argInfo{"--eval", argStr, argSetStr(&evalVal), " Give one line of source text"},
	argInfo{"--muffler", argStr, argSetStr(&mufflerVal), " Specify muffler color"},
}

func readArg() {
	argParse(argSpecList, func(arg string) error {
		inFile = arg
		return nil
	})
	if inFile == "" && evalVal == "" {
		scePanic(errors.New("no input file designation."))
	} else if inFile == "" && evalVal != "" {
		inFile = "(eval)"
	} else if inFile != "" && evalVal != "" {
		scePanic(errors.New("both input file and --eval are given."))
	}
	if outFile == "" {
		if evalVal == "" {
			outFile = changeExt(inFile, ".pdf")
		} else {
			outFile = "output.pdf"
		}
	}
	if mufflerVal == "" {
		mufflerVal = dfltMuffler
	}
	if textModeVal != "" {
		vals := strings.Split(textModeVal, ",")
		if len(vals) > 1 {
			scePanic(errors.New("--text-mode can have only one value."))
		}
		textMode = strings.TrimSpace(vals[0])
	}
	if markdownVal != "" {
		isMarkdown = true
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

	if isShowFont {
		fmt.Printf("  all the available fonts:")
		fmt.Printf("  ...oops, there's no such inessential concept as font!")
	}

	fmt.Printf("  evaluation done.\n")

	if byteComp {
		byteExec(byteCompile(value))
	} else {
		writeOutput(outFile)
	}
}

func writeOutput(pdst string) {
	switch textMode {
	case "":
		writePdf(pdst)
	case "plain":
		writeText(pdst, makePlainText())
	case "html":
		writeText(pdst, makeHtmlText())
	case "xml":
		writeText(pdst, makeXmlText())
	default:
		scePanic(fmt.Errorf("unknown text mode value '%s'", textMode))
	}

	fmt.Printf(" ---- ---- ---- ----\n")
	fmt.Printf("  output written on '%s'.\n", ordPath(pdst))
}

func writeText(pdst, text string) {
	wdst, err := os.Create(pdst)
	sceAssert(err)
	defer wdst.Close()

	_, err = io.WriteString(wdst, text)
	sceAssert(err)
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
}

func readFile(psrc string, value scValue) {
	fmt.Printf(" ---- ---- ---- ----\n")
	fmt.Printf("  reading '%s' ...\n", ordInPath(psrc))
	fmt.Printf("  type check passed. (%s)\n", value.vtype)

	if !typeCheckOnly && value.vtype != scEssential {
		scePanic(sceNonDocError(psrc, value.vtype))
	}
}

func openInFile(psrc string) (io.ReadCloser, error) {
	if evalVal == "" {
		return os.Open(psrc)
	}
	buf := bytes.NewBuffer([]byte(evalVal))
	return ioutil.NopCloser(buf), nil
}

func parseFile(psrc string) (value scValue) {
	fmt.Printf("  parsing '%s' ...\n", ordInPath(psrc))

	rsrc, err := openInFile(psrc)
	sceAssert(err)
	defer rsrc.Close()
	if isMarkdown {
		value, err = scParseMarkdown(rsrc)
	} else {
		value, err = scParseReader(rsrc)
	}
	sceAssert(err)
	return
}
