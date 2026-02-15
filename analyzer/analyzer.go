package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var logMethodNames = map[string]struct{}{
	"Debug": {}, "Info": {}, "Warn": {}, "Error": {}, "Fatal": {}, "Panic": {},
	"Debugw": {}, "Infow": {}, "Warnw": {}, "Errorw": {}, "Fatalw": {}, "Panicw": {},
	"Debugf": {}, "Infof": {}, "Warnf": {}, "Errorf": {}, "Fatalf": {}, "Panicf": {},
}

var Analyzer = &analysis.Analyzer{
	Name: "loglint",
	Doc:  "analyzer",
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

		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			name := sel.Sel.Name
			if _, ok := logMethodNames[name]; ok {
				pass.Reportf(sel.Sel.Pos(), "found candidate log call: %s", name)
			}
			return
		}

		if ident, ok := call.Fun.(*ast.Ident); ok {
			name := ident.Name
			if _, ok := logMethodNames[name]; ok {
				pass.Reportf(ident.Pos(), " found candidate log call: %s", name)
			}
			return
		}
	})

	return nil, nil
}
