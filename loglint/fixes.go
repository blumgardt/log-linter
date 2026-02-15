package loglint

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func fixLowercaseStart(msg string) string {
	i := 0
	for i < len(msg) {
		r, size := utf8.DecodeRuneInString(msg[i:])
		if r == utf8.RuneError && size == 0 {
			return msg
		}
		if unicode.IsSpace(r) {
			i += size
			continue
		}
		lr := unicode.ToLower(r)
		if lr == r {
			return msg
		}
		return msg[:i] + string(lr) + msg[i+size:]
	}
	return msg
}

func fixNoSpecialChars(msg string, extraAllowed string) string {
	var b strings.Builder
	b.Grow(len(msg))

	for _, r := range msg {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ':
			b.WriteRune(r)
		case r == '_' || r == '-':
			b.WriteRune(r)
		default:
			if extraAllowed != "" && strings.ContainsRune(extraAllowed, r) {
				b.WriteRune(r)
				continue
			}
			if unicode.IsSpace(r) {
				b.WriteRune(' ')
			} else {
				b.WriteRune('_')
			}
		}
	}

	return b.String()
}
