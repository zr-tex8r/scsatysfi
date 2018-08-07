// Copyright (c) 2018 Takayuki YATO (aka. "ZR")
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
	} else {
		return filepath.Base(path)
	}
}

func changeExt(path, ext string) string {
	s := path[:len(path)-len(filepath.Ext(path))]
	return s + ext
}
