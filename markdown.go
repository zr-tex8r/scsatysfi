// Copyright (c) 2018 Takayuki YATO (aka. "ZR")
//   GitHub:   https://github.com/zr-tex8r
//   Twitter:  @zr_tex8r
// Distributed under the MIT License.

package main

import (
	"io"
)

func scParseMarkdown(rsrc io.Reader) (value scValue, err error) {
	if _, err = io.Copy(&nullWriter{}, rsrc); err != nil {
		return
	}
	// there is no such thing as "invalid" Markdown
	value = scValue{scEssential}
	return
}
