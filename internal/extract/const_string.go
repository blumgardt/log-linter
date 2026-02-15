package extract

import (
	"go/ast"
	"go/constant"
	"go/token"
	"strconv"

	"golang.org/x/tools/go/analysis"
)

func ExtractConstMessageAt(pass *analysis.Pass, call *ast.CallExpr, msgIndex int) (string, ast.Expr, bool) {
	if msgIndex < 0 || msgIndex >= len(call.Args) {
		return "", nil, false
	}
	expr := call.Args[msgIndex]
	s, ok := ExtractConstString(pass, expr)
	if !ok {
		return "", nil, false
	}
	return s, expr, true
}

func ExtractConstString(pass *analysis.Pass, expr ast.Expr) (string, bool) {
	if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		s, err := strconv.Unquote(lit.Value)
		if err != nil {
			return "", false
		}
		return s, true
	}

	if pass.TypesInfo != nil {
		if tv, ok := pass.TypesInfo.Types[expr]; ok && tv.Value != nil && tv.Value.Kind() == constant.String {
			return constant.StringVal(tv.Value), true
		}
	}

	return "", false
}
