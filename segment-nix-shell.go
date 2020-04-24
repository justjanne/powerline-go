package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
)

func segmentNixShell(p *powerline) {
	var nixShell string
	nixShell, _ = os.LookupEnv("IN_NIX_SHELL")
	if nixShell == "" {
		return
	}
	p.appendSegment("nix-shell", pwl.Segment{
		Content:    nixShell,
		Foreground: p.theme.NixShellFg,
		Background: p.theme.NixShellBg,
	})
}
