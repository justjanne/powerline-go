package main

import (
	"os"
)

func segmentNixShell(p *powerline) {
	var nixShell string
	nixShell, _ = os.LookupEnv("IN_NIX_SHELL")
	if nixShell == "" {
		return
	}
	p.appendSegment("nix-shell", segment{
		content:    nixShell,
		foreground: p.theme.NixShellFg,
		background: p.theme.NixShellBg,
	})
}
