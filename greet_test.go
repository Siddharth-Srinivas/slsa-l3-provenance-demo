package main

import (
	"strings"
	"testing"
)

func TestGreeting(t *testing.T) {
	tests := []struct {
		name, in, want string
	}{
		{"plain name", "Ada", "Hello, Ada!"},
		{"empty falls back to world", "", "Hello, world!"},
		{"surrounding space is trimmed", "  Ada  ", "Hello, Ada!"},
		{"control characters are dropped", "Ada\x1b[2J\r\n", "Hello, Ada[2J!"},
		{"bidi overrides are dropped", "Ada\u202egnp.exe", "Hello, Adagnp.exe!"},
		{"only unprintable falls back to world", "\x00\x07\u200b", "Hello, world!"},
		{"long names are truncated", strings.Repeat("é", maxNameLen+5), "Hello, " + strings.Repeat("é", maxNameLen) + "!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Greeting(tt.in); got != tt.want {
				t.Errorf("Greeting(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
