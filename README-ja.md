![logo1](https://raw.githubusercontent.com/zr-tex8r/scsatysfi/images/scsatysfi-logo.png)

![Snowman Status](https://raw.githubusercontent.com/zr-tex8r/scsatysfi/images/snowman-nice-green.png)

## 概要

*scSATySFi*（テキトーに発音します）は、画期的な組版処理システムとその言語です。構文は主に☃部分と非☃部分からなり、前者は本質的な内容を記述し、後者は本質的でないので無視されます。いわゆる“SC版”のソフトウェアのため、常に本質的な出力が実現されています。

本ソフトウェアは2018年度ﾅﾝﾄｶの日の1ネタとして開発されました。

## インストール方法

フツーにGo言語の処理系をインストールします。

※ `$GOPATH/bin`を実行パスに追加してください。

その後、以下のコマンドを実行すると`scsatysfi`の実行ファイルがインストールされます：

    go get github.com/zr-tex8r/scsatysfi

## 用法

    scsatysfi <input files> -o <output file>

で`<input files>`から`<output file>`を出力します。例えばソースファイル`duck.scty`から`essential.pdf`を出力したい場合、次のようにします：

    scsatysfi duck.scty -o essential.pdf

## コマンドラインオプション

  * `-v`／`--version`：バージョンを表示します。
  * `-o`／`--output`：出力ファイル名を指定します。省略された場合、入力ファイル名の拡張子を`.pdf`に変えた名前を出力ファイル名とします。
  * `--full-path`：標準出力に書き込むログに於いて、ファイル名をすべて絶対パスで表示します。
  * `--type-check-only`：型検査だけをして終了します。

## ライセンス

MITライセンスが適用されます。

--------------------
Takayuki YATO (aka. "ZR")  
https://github.com/zr-tex8r
