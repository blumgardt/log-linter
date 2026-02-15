package rules_test

import (
	"testing"

	"github.com/blumgardt/log-linter/internal/config"
	"github.com/blumgardt/log-linter/internal/rules"
)

func TestViolatesNoSpecialChars_DefaultAllowlist(t *testing.T) {
	config.ExtraAllowedChars = ""

	tests := []struct {
		name string
		msg  string
		want bool
	}{
		{"ok_simple", "server started", false},
		{"ok_dash_underscore", "server_started-ok", false},
		{"ok_digits", "id 123", false},

		{"bang", "server started!", true},
		{"dot", "server.started", true},
		{"percent", "user=%s", true},
		{"slash", "path /home", true},
		{"colon", "x:y", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rules.ViolatesNoSpecialChars(tt.msg)
			if got != tt.want {
				t.Fatalf("ViolatesNoSpecialChars(%q)=%v, want %v", tt.msg, got, tt.want)
			}
		})
	}
}

func TestViolatesNoSpecialChars_ExtraAllowed(t *testing.T) {
	config.ExtraAllowedChars = ":/."
	defer func() { config.ExtraAllowedChars = "" }()

	tests := []struct {
		name string
		msg  string
		want bool
	}{
		{"colon_allowed", "x:y", false},
		{"slash_allowed", "path /home", false},
		{"dot_allowed", "server.started", false},
		{"bang_still_forbidden", "server!", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rules.ViolatesNoSpecialChars(tt.msg)
			if got != tt.want {
				t.Fatalf("ViolatesNoSpecialChars(%q)=%v, want %v", tt.msg, got, tt.want)
			}
		})
	}
}
