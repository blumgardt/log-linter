package rules

import (
	"go/ast"
	"strings"

	"github.com/blumgardt/log-linter/internal/config"
	"github.com/blumgardt/log-linter/internal/detect"
	"github.com/blumgardt/log-linter/internal/extract"

	"golang.org/x/tools/go/analysis"
)

func findSensitiveKeyword(s string) (string, bool) {
	low := strings.ToLower(s)
	for _, kw := range config.SensitiveKeywords {
		if strings.Contains(low, kw) {
			return kw, true
		}
	}
	return "", false
}

func ReportSensitive(pass *analysis.Pass, call *ast.CallExpr, msgExpr ast.Expr, msg string, info detect.CallInfo) bool {
	if hit, ok := findSensitiveKeyword(msg); ok {
		pass.Reportf(msgExpr.Pos(), "log message contains sensitive keyword %q", hit)
		return true
	}

	argsStart := info.MsgIndex + 1

	// slog / sugared zap
	if info.Kind == detect.LoggerSlog || (info.Kind == detect.LoggerZapSugared && strings.HasSuffix(info.Method, "w")) {
		for i := argsStart; i+1 < len(call.Args); i += 2 {
			keyExpr := call.Args[i]
			key, ok := extract.ExtractConstString(pass, keyExpr)
			if !ok {
				continue
			}
			if hit, ok := findSensitiveKeyword(key); ok {
				pass.Reportf(keyExpr.Pos(), "log field key contains sensitive keyword %q", hit)
				return true
			}

			valExpr := call.Args[i+1]
			if ident, ok := valExpr.(*ast.Ident); ok {
				if hit, ok := findSensitiveKeyword(ident.Name); ok {
					pass.Reportf(ident.Pos(), "logging potentially sensitive value %q", hit)
					return true
				}
			}
		}
	}

	// zap fields
	for _, arg := range call.Args[argsStart:] {
		fieldCall, ok := arg.(*ast.CallExpr)
		if !ok {
			continue
		}
		if _, ok := fieldCall.Fun.(*ast.SelectorExpr); !ok {
			continue
		}

		if len(fieldCall.Args) >= 1 {
			keyExpr := fieldCall.Args[0]
			key, ok := extract.ExtractConstString(pass, keyExpr)
			if ok {
				if hit, ok := findSensitiveKeyword(key); ok {
					pass.Reportf(keyExpr.Pos(), "zap field key contains sensitive keyword %q", hit)
					return true
				}
			}
		}

		if len(fieldCall.Args) >= 2 {
			if ident, ok := fieldCall.Args[1].(*ast.Ident); ok {
				if hit, ok := findSensitiveKeyword(ident.Name); ok {
					pass.Reportf(ident.Pos(), "logging potentially sensitive value %q", hit)
					return true
				}
			}
		}
	}

	return false
}
