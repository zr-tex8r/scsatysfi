// Copyright (c) 2018 Takayuki YATO (aka. "ZR")
//   GitHub:   https://github.com/zr-tex8r
//   Twitter:  @zr_tex8r
// Distributed under the MIT License.

package main

import (
	"fmt"
	"os"
)

func init() {
	initSceUnxDesc()
}

type sceError struct {
	tag     string
	message string
}

const (
	sceSynTag  = "Syntax Error at Lexer"
	sceTypeTag = "Type Error"
	sceMiscTag = "Error"
)

func (e *sceError) Error() string {
	return e.tag + ": " + e.message
}

func sceDesc(e error) string {
	tag, msg := sceMiscTag, ""
	switch e := e.(type) {
	case *sceError:
		tag, msg = e.tag, e.message
	case *os.PathError:
		msg = natFullPath(e.Path) + ": " + sceUnxDesc(e.Err)
	default:
		msg = e.Error()
	}
	return fmt.Sprintf("! [%v] %v\n", tag, msg)
}

func scePanic(err error) {
	fmt.Print(sceDesc(err))
	os.Exit(1)
}

func sceAssert(err error) {
	if err != nil {
		scePanic(err)
	}
}

func scePosDesc(line, bcol, ecol int, msg string) string {
	return fmt.Sprintf(
		"at line %v, characters %v-%v:\n    %v",
		line, bcol, ecol, msg)
}

func sceBadCharError(line, bcol, ecol int, chr rune) error {
	msg := fmt.Sprintf("invalid character %q(%U)", chr, chr)
	return &sceError{sceSynTag, scePosDesc(line, bcol, ecol, msg)}
}

func sceBadCommentError(line int) error {
	msg := fmt.Sprintf("text input ended while reading a block comment")
	return &sceError{sceSynTag, scePosDesc(line, 0, 0, msg)}
}

func sceNonDocError(path string, vt scVType) error {
	msg := fmt.Sprintf(
		"file '%v' is not an essential file; it is of type\n      %v",
		fullInPath(path), vt)
	return &sceError{sceTypeTag, msg}
}

//-------- sceUnxDesc

var sceErrDescMap map[string]string

func initSceUnxDesc() {
	sceErrDescMap = make(map[string]string, len(sceErrCodeMap))
	for wc, uc := range sceErrCodeMap {
		sceErrDescMap[sceWinErrDesc[wc]] = sceUnxErrDesc[uc]
	}
}

func sceUnxDesc(e error) string {
	desc := e.Error()
	if udesc, ok := sceErrDescMap[desc]; ok {
		return udesc
	}
	return desc
}

var sceErrCodeMap = map[int]int{
	1:    22,
	2:    2,
	3:    2,
	4:    24,
	5:    13,
	6:    9,
	7:    12,
	8:    12,
	9:    12,
	12:   22,
	13:   22,
	15:   2,
	32:   13,
	33:   13,
	80:   17,
	82:   13,
	87:   22,
	108:  13,
	109:  32,
	112:  28,
	114:  9,
	131:  22,
	145:  41,
	158:  13,
	161:  2,
	164:  11,
	167:  13,
	183:  17,
	206:  2,
	215:  11,
	1816: 12,
}

var sceWinErrDesc = map[int]string{
	1:    "Incorrect function.",
	2:    "The system cannot find the file specified.",
	3:    "The system cannot find the path specified.",
	4:    "The system cannot open the file.",
	5:    "Access is denied.",
	6:    "The handle is invalid.",
	7:    "The storage control blocks were destroyed.",
	8:    "Not enough storage is available to process this command.",
	9:    "The storage control block address is invalid.",
	12:   "The access code is invalid.",
	13:   "The data is invalid.",
	15:   "The system cannot find the drive specified.",
	32:   "The process cannot access the file because it is being used by another process.",
	33:   "The process cannot access the file because another process has locked a portion of the file.",
	80:   "The file exists.",
	82:   "The directory or file cannot be created.",
	87:   "The parameter is incorrect.",
	108:  "The disk is in use or locked by another process.",
	109:  "The pipe has been ended.",
	112:  "There is not enough space on the disk.",
	114:  "The target internal file identifier is incorrect.",
	131:  "An attempt was made to move the file pointer before the beginning of the file.",
	145:  "The directory is not empty.",
	158:  "The segment is already unlocked.",
	161:  "The specified path is invalid.",
	164:  "No more threads can be created in the system.",
	167:  "Unable to lock a region of a file.",
	183:  "Cannot create a file when that file already exists.",
	206:  "The filename or extension is too long.",
	215:  "Cannot nest calls to LoadModule.",
	1816: "Not enough quota is available to process this command.",
}

var sceUnxErrDesc = map[int]string{
	2:  "No such file or directory",
	9:  "Bad file descriptor",
	11: "Resource temporarily unavailable",
	12: "Not enough space",
	13: "Permission denied",
	17: "File exists",
	22: "Invalid argument",
	24: "Too many open files",
	28: "No space left on device",
	32: "Broken pipe",
	41: "Directory not empty",
}
