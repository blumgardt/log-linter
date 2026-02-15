package loglint

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/blumgardt/log-linter/internal/config"
	"github.com/blumgardt/log-linter/internal/detect"
	"github.com/blumgardt/log-linter/internal/extract"
	"github.com/blumgardt/log-linter/internal/rules"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "loglint",
	Doc:  "loglint",
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	Run: run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	ins.Preorder(nodeFilter, func(node ast.Node) {
		call := node.(*ast.CallExpr)

		methodName, ok := detect.ExtractCalledName(call.Fun)
		if !ok {
			return
		}
		if !detect.IsLogMethodName(methodName) {
			return
		}

		info, ok := detect.DetectSupportedLogCall(pass, call)
		if !ok {
			return
		}

		msg, msgExpr, ok := extract.ExtractConstMessageAt(pass, call, info.MsgIndex)
		if !ok {
			return
		}

		if config.IsRule4Enabled() {
			if rules.ReportSensitive(pass, call, msgExpr, msg, info) {
				return
			}
		}

		if config.IsRule1Enabled() && rules.ViolatesLowercaseStart(msg) {
			fixed := fixLowercaseStart(msg)
			reportWithMsgFix(pass, msgExpr, "log message must start with a lower case letter", fixed)
		}

		if config.IsRule2Enabled() && rules.ViolatesEnglishOnlyASCII(msg) {
			pass.Reportf(msgExpr.Pos(), "log message must contains only English characters")
			return
		}

		if config.IsRule3Enabled() && rules.ViolatesNoSpecialChars(msg) {
			fixed := fixNoSpecialChars(msg, config.GetExtraAllowedChars())
			reportWithMsgFix(pass, msgExpr, "log message must not contain any special characters or emoji", fixed)
		}
	})

	return nil, nil
}

func reportWithMsgFix(pass *analysis.Pass, msgExpr ast.Expr, message string, fixedMsg string) {
	diag := analysis.Diagnostic{
		Pos:     msgExpr.Pos(),
		Message: message,
	}

	lit, ok := msgExpr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		pass.Report(diag)
		return
	}

	orig, err := strconv.Unquote(lit.Value)
	if err != nil {
		pass.Report(diag)
		return
	}

	if fixedMsg == "" || fixedMsg == orig {
		pass.Report(diag)
		return
	}

	newLit := strconv.Quote(fixedMsg)

	diag.SuggestedFixes = []analysis.SuggestedFix{
		{
			Message: "apply suggested fix",
			TextEdits: []analysis.TextEdit{
				{
					Pos:     lit.Pos(),
					End:     lit.End(),
					NewText: []byte(newLit),
				},
			},
		},
	}

	pass.Report(diag)
}
