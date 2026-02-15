package main

import (
	loglinter "github.com/blumgardt/log-linter/loglint"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(loglinter.Analyzer)
}
