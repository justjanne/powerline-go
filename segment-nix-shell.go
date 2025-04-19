package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
)

func segmentNixShell(p *powerline) []pwl.Segment {
	var nixShell string
	nixShell, _ = os.LookupEnv("IN_NIX_SHELL")
	if nixShell == "" {
		return []pwl.Segment{}
	}
	return []pwl.Segment{{
		Name:       "nix-shell",
		Content:    p.symbols.NixShellIndicator,
		Foreground: p.theme.NixShellFg,
		Background: p.theme.NixShellBg,
	}}
}
