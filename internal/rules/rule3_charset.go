package rules

import (
	"github.com/blumgardt/log-linter/internal/config"
)

func ViolatesNoSpecialChars(msg string) bool {
	extra := config.GetExtraAllowedChars()

	for _, r := range msg {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == ' ':
		case r == '_' || r == '-':
		default:
			if extra != "" && containsRune(extra, r) {
				continue
			}
			return true
		}
	}
	return false
}

func containsRune(s string, target rune) bool {
	for _, r := range s {
		if r == target {
			return true
		}
	}
	return false
}
