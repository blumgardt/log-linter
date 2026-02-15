package rules

import (
	"strings"
	"unicode"
)

func ViolatesLowercaseStart(msg string) bool {
	s := strings.TrimLeftFunc(msg, unicode.IsSpace)
	if s == "" {
		return false
	}
	r, _ := utf8FirstRune(s)
	return unicode.IsUpper(r)
}

func utf8FirstRune(s string) (rune, int) {
	for _, r := range s {
		return r, len(string(r))
	}
	return 0, 0
}
