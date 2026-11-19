package main

import (
	"strings"
	"unicode"
)

// maxNameLen caps how many runes of the name are echoed back.
const maxNameLen = 64

// Greeting returns a greeting for name. The name is untrusted input, so runes
// that aren't printable are dropped: control characters (including the ESC
// that starts terminal escape sequences), bidi overrides and other invisible
// formatting characters. The result is truncated to maxNameLen runes.
func Greeting(name string) string {
	var b strings.Builder
	n := 0
	for _, r := range strings.TrimSpace(name) {
		if n == maxNameLen {
			break
		}
		if unicode.IsPrint(r) {
			b.WriteRune(r)
			n++
		}
	}
	clean := strings.TrimSpace(b.String())
	if clean == "" {
		clean = "world"
	}
	return "Hello, " + clean + "!"
}
