package main

import (
	"fmt"
	"strings"
	"testing"

	pwl "github.com/justjanne/powerline-go/powerline"
)

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

func newTestPowerline(shell string, align alignment, rows int) *powerline {
	p := &powerline{
		cfg:     defaults,
		shell:   defaults.Shells[shell],
		theme:   defaults.Themes["default"],
		symbols: defaults.Modes["patched"],
		align:   align,
	}
	p.reset = fmt.Sprintf(p.shell.ColorTemplate, "[0m")
	for i := 0; i < rows; i++ {
		p.Segments = append(p.Segments, []pwl.Segment{{
			Name:                "exit",
			Content:             "SIGPIPE",
			Foreground:          p.theme.CmdFailedFg,
			Background:          p.theme.CmdFailedBg,
			Separator:           p.symbols.Separator,
			SeparatorForeground: p.theme.CmdFailedBg,
		}})
	}
	return p
}

// A row that soft-wraps on the last screen line makes the terminal scroll in a
// line filled with the background colour active at that moment, so rows that
// are followed by another line end by erasing to the end of the line. The line
// the user types on must never contain that erase: readline redraws it from
// column 0 on every keystroke and would wipe the typed command.
func Test_drawErasesToEndOfLineOnlyAboveInputLine(t *testing.T) {
	tests := []struct {
		name    string
		shell   string
		align   alignment
		rows    int
		newline bool
		want    int // number of erase sequences, each directly followed by a line break
	}{
		{"bash single row", "bash", alignLeft, 1, false, 0},
		{"bash single row with -newline", "bash", alignLeft, 1, true, 1},
		{"bash two rows", "bash", alignLeft, 2, false, 1},
		{"zsh two rows with -newline", "zsh", alignLeft, 2, true, 2},
		{"bare two rows", "bare", alignLeft, 2, false, 1},
		{"zsh right prompt", "zsh", alignRight, 1, true, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newTestPowerline(tt.shell, tt.align, tt.rows)
			p.cfg.PromptOnNewLine = tt.newline
			got := p.draw()
			erase := fmt.Sprintf(p.shell.ColorTemplate, "[K")
			if n := strings.Count(got, erase); n != tt.want {
				t.Errorf("draw() = %q, contains %d erase sequences, want %d", got, n, tt.want)
			}
			if n := strings.Count(got, erase+"\n"); n != tt.want {
				t.Errorf("draw() = %q, every erase sequence must directly precede a line break", got)
			}
		})
	}
}
