package rules_test

import (
	"go/ast"
	"go/token"
	"strings"
	"testing"

	"github.com/blumgardt/log-linter/internal/config"
	"github.com/blumgardt/log-linter/internal/detect"
	"github.com/blumgardt/log-linter/internal/rules"
	"golang.org/x/tools/go/analysis"
)

func newPassCapture(diags *[]analysis.Diagnostic) *analysis.Pass {
	return &analysis.Pass{
		Report: func(d analysis.Diagnostic) {
			*diags = append(*diags, d)
		},
	}
}

func TestReportSensitive_MessageKeyword(t *testing.T) {
	config.SensitiveKeywords = []string{"token"}

	var diags []analysis.Diagnostic
	pass := newPassCapture(&diags)

	msgExpr := &ast.BasicLit{Kind: token.STRING, Value: `"token leaked"`}
	call := &ast.CallExpr{
		Args: []ast.Expr{msgExpr},
	}
	info := detect.CallInfo{Kind: detect.LoggerSlog, Method: "Info", MsgIndex: 0}

	ok := rules.ReportSensitive(pass, call, msgExpr, "token leaked", info)
	if !ok {
		t.Fatalf("expected ReportSensitive=true")
	}
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
	if !strings.Contains(diags[0].Message, `sensitive keyword "token"`) {
		t.Fatalf("unexpected diagnostic message: %q", diags[0].Message)
	}
}

func TestReportSensitive_StructuredKey(t *testing.T) {
	config.SensitiveKeywords = []string{"password"}

	var diags []analysis.Diagnostic
	pass := newPassCapture(&diags)

	msgExpr := &ast.BasicLit{Kind: token.STRING, Value: `"login"`}
	keyExpr := &ast.BasicLit{Kind: token.STRING, Value: `"password"`}
	valExpr := &ast.Ident{Name: "password"}

	call := &ast.CallExpr{
		Args: []ast.Expr{msgExpr, keyExpr, valExpr},
	}
	info := detect.CallInfo{Kind: detect.LoggerSlog, Method: "Info", MsgIndex: 0}

	ok := rules.ReportSensitive(pass, call, msgExpr, "login", info)
	if !ok {
		t.Fatalf("expected ReportSensitive=true")
	}
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
	if !strings.Contains(diags[0].Message, `log field key contains sensitive keyword "password"`) {
		t.Fatalf("unexpected diagnostic message: %q", diags[0].Message)
	}
}

func TestReportSensitive_StructuredValueIdent(t *testing.T) {
	config.SensitiveKeywords = []string{"token"}

	var diags []analysis.Diagnostic
	pass := newPassCapture(&diags)

	msgExpr := &ast.BasicLit{Kind: token.STRING, Value: `"login"`}
	keyExpr := &ast.BasicLit{Kind: token.STRING, Value: `"id"`}
	valExpr := &ast.Ident{Name: "token"}

	call := &ast.CallExpr{
		Args: []ast.Expr{msgExpr, keyExpr, valExpr},
	}
	info := detect.CallInfo{Kind: detect.LoggerSlog, Method: "Info", MsgIndex: 0}

	ok := rules.ReportSensitive(pass, call, msgExpr, "login", info)
	if !ok {
		t.Fatalf("expected ReportSensitive=true")
	}
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
	if !strings.Contains(diags[0].Message, `logging potentially sensitive value "token"`) {
		t.Fatalf("unexpected diagnostic message: %q", diags[0].Message)
	}
}

func TestReportSensitive_ZapFieldKey(t *testing.T) {
	config.SensitiveKeywords = []string{"token"}

	var diags []analysis.Diagnostic
	pass := newPassCapture(&diags)

	msgExpr := &ast.BasicLit{Kind: token.STRING, Value: `"login"`}
	fieldCall := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "zap"},
			Sel: &ast.Ident{Name: "String"},
		},
		Args: []ast.Expr{
			&ast.BasicLit{Kind: token.STRING, Value: `"token"`},
			&ast.Ident{Name: "token"},
		},
	}

	call := &ast.CallExpr{
		Args: []ast.Expr{msgExpr, fieldCall},
	}
	info := detect.CallInfo{Kind: detect.LoggerZap, Method: "Info", MsgIndex: 0}

	ok := rules.ReportSensitive(pass, call, msgExpr, "login", info)
	if !ok {
		t.Fatalf("expected ReportSensitive=true")
	}
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
	if !strings.Contains(diags[0].Message, `zap field key contains sensitive keyword "token"`) {
		t.Fatalf("unexpected diagnostic message: %q", diags[0].Message)
	}
}
