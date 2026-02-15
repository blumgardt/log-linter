package loglint_test

import (
	"testing"

	loglint "github.com/blumgardt/log-linter/loglint"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, loglint.Analyzer, "a")
}
