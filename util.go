// Copyright (c) 2018-2021 Takayuki YATO (aka. "ZR")
//   GitHub:   https://github.com/zr-tex8r
//   Twitter:  @zr_tex8r
// Distributed under the MIT License.

package main

import (
	"path/filepath"
)

func unxFullPath(path string) string {
	r, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(r)
}

func natFullPath(path string) string {
	if r, err := filepath.Abs(path); err == nil {
		return r
	}
	return path
}

func ordPath(path string) string {
	if fullPath {
		return natFullPath(path)
	}
	return filepath.Base(path)
}

func ordInPath(path string) string {
	if fullPath && evalVal == "" {
		return natFullPath(path)
	}
	return filepath.Base(path)
}

func fullInPath(path string) string {
	if evalVal == "" {
		return natFullPath(path)
	}
	return path
}

func changeExt(path, ext string) string {
	s := path[:len(path)-len(filepath.Ext(path))]
	return s + ext
}

//--------nullWriter

type nullWriter struct{}

func (w *nullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
