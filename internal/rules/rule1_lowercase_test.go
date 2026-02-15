package rules_test

import (
	"testing"

	"github.com/blumgardt/log-linter/internal/rules"
)

func TestViolatesLowercaseStart(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want bool
	}{
		{"empty", "", false},
		{"spaces_only", "   \t\n", false},
		{"lowercase", "server started", false},
		{"uppercase", "Server started", true},
		{"leading_spaces_uppercase", "   Server started", true},
		{"leading_spaces_lowercase", "   server started", false},
		{"digit_first", "1 server started", false},
		{"underscore_first", "_server started", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rules.ViolatesLowercaseStart(tt.msg)
			if got != tt.want {
				t.Fatalf("ViolatesLowercaseStart(%q)=%v, want %v", tt.msg, got, tt.want)
			}
		})
	}
}
