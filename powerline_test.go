package main

import "testing"

func Test_detectShell(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "",
			want: "bare",
		},
		{
			name: "/bin/sh",
			want: "bare",
		}, {
			name: "/bin/bash",
			want: "bash",
		}, {
			name: "/usr/local/bin/bash5",
			want: "bash",
		}, {
			name: "/usr/bin/zsh",
			want: "zsh",
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := detectShell(tt.name); got != tt.want {
				t.Errorf("detectShell(%q) = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}
