package util

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func CamelCase(s string) string {
	var parts = strings.Split(s, "_")
	for i, part := range parts {
		parts[i] = Ucfirst(part)
	}
	return strings.Join(parts, "")
}

func Ucfirst(s string) string {
	var b = []byte(s)
	r, _ := utf8.DecodeRuneInString(s)
	r = unicode.ToUpper(r)
	utf8.EncodeRune(b, r)
	return string(b)
}
