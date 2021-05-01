// Copyright (c) 2018-2021 Takayuki YATO (aka. "ZR")
//   GitHub:   https://github.com/zr-tex8r
//   Twitter:  @zr_tex8r
// Distributed under the MIT License.

package main

import (
	_ "errors"
	_ "fmt"
	"github.com/zr-tex8r/xcolor"
)

func sctHtmlMufflerColor() string {
	if mufflerVal == dfltMuffler {
		return ""
	}

	col, err := xcolor.Parse(mufflerVal)
	if err != nil {
		scePanic(err)
	}
	return col.HtmlCode()
}

//-------- plain text

const sctSnowman = "     " + // HURRAY!
	`  ____         
    ___HHHH   _____ 
   / .   . \ |NICE!|
   \  ---  / |~~~~~ 
 V :#######: Y      
  \/   o*"*\/       
  {    o    }       
   \_______/        
`

func makePlainText() string {
	return sctSnowman
}

//-------- HTML

const sctHtmlPrologue = // something HTML5
`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
</head>
<style>
body {
font-size: 480px;
font-family: "SCAlleSnowman", sans-serif;
}
</style>
<body>
`

const sctHtmlEpilogue = //
`</body>
</html>
`

func makeHtmlText() string {
	snowman := "&#x2603;\n"
	return sctHtmlPrologue + snowman + sctHtmlEpilogue
}

//-------- XML

func makeXmlText() string {
	col := sctHtmlMufflerColor()
	if col != "" {
		col = "muffler=\"" + col + "\""
	}
	return "<snowman " + col + "/>\n"
}
