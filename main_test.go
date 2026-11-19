package main

import (
	"bytes"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantCode int
		wantOut  string
	}{
		{"default greeting", nil, 0, "Hello, world!\n"},
		{"name flag", []string{"--name", "Ada"}, 0, "Hello, Ada!\n"},
		{"version flag", []string{"--version"}, 0, "dev\n"},
		{"help", []string{"-h"}, 0, ""},
		{"unknown flag", []string{"--nope"}, 2, ""},
		{"positional argument", []string{"Ada"}, 2, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run(tt.args, &stdout, &stderr)
			if code != tt.wantCode || stdout.String() != tt.wantOut {
				t.Errorf("run(%q) = %d, %q; want %d, %q (stderr: %q)",
					tt.args, code, stdout.String(), tt.wantCode, tt.wantOut, stderr.String())
			}
		})
	}
}
