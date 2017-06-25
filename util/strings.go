package util

import (
	"strings"
	"unicode"
)

func StripNonSpaceWhitespaces(input string) string {
	return strings.Map(func(r rune) rune {
		if isSpace(r) {
			return -1
		}
		return r
	}, input)
}

func isSpace(r rune) bool {
	// This property isn't the same as Z; special-case it.
	if uint32(r) <= unicode.MaxLatin1 {
		switch r {
		case '\t', '\n', '\v', '\f', '\r', 0x85, 0xA0:
			return true
		}
		return false
	}
	return false
}
