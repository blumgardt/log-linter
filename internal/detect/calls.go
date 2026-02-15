package detect

import (
	"go/ast"

	"github.com/blumgardt/log-linter/internal/config"
)

func ExtractCalledName(fun ast.Expr) (string, bool) {
	if sel, ok := fun.(*ast.SelectorExpr); ok {
		return sel.Sel.Name, true
	}
	if ident, ok := fun.(*ast.Ident); ok {
		return ident.Name, true
	}
	return "", false
}

func IsLogMethodName(name string) bool {
	_, ok := config.LogMethodNames[name]
	return ok
}
