package detect

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

type LoggerKind int

const (
	LoggerUnknown LoggerKind = iota
	LoggerSlog
	LoggerZap
	LoggerZapSugared
)

type CallInfo struct {
	Kind     LoggerKind
	Method   string
	MsgIndex int
}

var slogNames = map[string]struct{}{
	"Debug": {}, "Info": {}, "Warn": {}, "Error": {},
	"DebugContext": {}, "InfoContext": {}, "WarnContext": {}, "ErrorContext": {},
}

var zapLoggerNames = map[string]struct{}{
	"Debug": {}, "Info": {}, "Warn": {}, "Error": {}, "DPanic": {}, "Panic": {}, "Fatal": {},
}

var zapSugaredNames = map[string]struct{}{
	"Debug": {}, "Info": {}, "Warn": {}, "Error": {}, "DPanic": {}, "Panic": {}, "Fatal": {},
	"Debugf": {}, "Infof": {}, "Warnf": {}, "Errorf": {}, "DPanicf": {}, "Panicf": {}, "Fatalf": {},
	"Debugw": {}, "Infow": {}, "Warnw": {}, "Errorw": {}, "DPanicw": {}, "Panicw": {}, "Fatalw": {},
}

func DetectSupportedLogCall(pass *analysis.Pass, call *ast.CallExpr) (CallInfo, bool) {
	switch fun := call.Fun.(type) {

	case *ast.SelectorExpr:
		if sel := pass.TypesInfo.Selections[fun]; sel != nil {
			meth := sel.Obj()
			if meth == nil || meth.Pkg() == nil {
				return CallInfo{}, false
			}
			pkgPath := meth.Pkg().Path()
			name := meth.Name()

			recv := sel.Recv()
			recvPkgPath, recvTypeName := namedType(recv)

			// slog
			if pkgPath == "log/slog" && recvPkgPath == "log/slog" && recvTypeName == "Logger" {
				if _, ok := slogNames[name]; ok {
					return CallInfo{Kind: LoggerSlog, Method: name, MsgIndex: 0}, true
				}
			}

			// zap
			if pkgPath == "go.uber.org/zap" && recvPkgPath == "go.uber.org/zap" {
				if recvTypeName == "Logger" {
					if _, ok := zapLoggerNames[name]; ok {
						return CallInfo{Kind: LoggerZap, Method: name, MsgIndex: 0}, true
					}
				}
				if recvTypeName == "SugaredLogger" {
					if _, ok := zapSugaredNames[name]; ok {

						return CallInfo{Kind: LoggerZapSugared, Method: name, MsgIndex: 0}, true
					}
				}
			}
			return CallInfo{}, false
		}

		if obj := pass.TypesInfo.Uses[fun.Sel]; obj != nil && obj.Pkg() != nil {
			pkgPath := obj.Pkg().Path()
			name := obj.Name()

			if pkgPath == "log/slog" {
				if _, ok := slogNames[name]; ok {
					msgIndex := 0
					if strings.HasSuffix(name, "Context") {
						msgIndex = 1
					}
					return CallInfo{Kind: LoggerSlog, Method: name, MsgIndex: msgIndex}, true
				}
			}

			return CallInfo{}, false
		}

	case *ast.Ident:
		return CallInfo{}, false
	}

	return CallInfo{}, false
}

func namedType(t types.Type) (pkgPath, typeName string) {
	for {
		if p, ok := t.(*types.Pointer); ok {
			t = p.Elem()
			continue
		}
		break
	}
	n, ok := t.(*types.Named)
	if !ok {
		return "", ""
	}
	obj := n.Obj()
	if obj == nil || obj.Pkg() == nil {
		return "", ""
	}
	return obj.Pkg().Path(), obj.Name()
}
