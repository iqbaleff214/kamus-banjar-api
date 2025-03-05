package utils

import (
	"strings"
	"unicode"
)

func IsAlphabet(s string) bool {
	return !strings.ContainsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
}
