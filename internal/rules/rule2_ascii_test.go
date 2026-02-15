package rules_test

import (
	"testing"

	"github.com/blumgardt/log-linter/internal/rules"
)

func TestViolatesEnglishOnlyASCII(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want bool
	}{
		{"ascii_ok", "server started_123 - ok", false},
		{"cyrillic", "сервер стартанул", true},
		{"mixed", "server стартанул", true},
		{"emoji", "server 😀", true},
		{"accented", "café", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rules.ViolatesEnglishOnlyASCII(tt.msg)
			if got != tt.want {
				t.Fatalf("ViolatesEnglishOnlyASCII(%q)=%v, want %v", tt.msg, got, tt.want)
			}
		})
	}
}
