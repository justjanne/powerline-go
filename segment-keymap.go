package main

import (
	"os"
)

func segmentKeymap(p *powerline) {
	var keymap string
	var present bool
	if keymap == "" {
		keymap, present = os.LookupEnv("KEYMAP_POWERLINE")
	}
	if !present {
		keymap = "ouch"
	}
	if keymap == "main" {
		p.appendSegment("keymap", segment{
			content:    " I ",
			foreground: p.theme.RepoDirtyFg,
			background: p.theme.RepoDirtyBg,
		})
	}
	if keymap == "vicmd" {
		p.appendSegment("keymap", segment{
			content:    " C ",
			foreground: p.theme.RepoCleanFg,
			background: p.theme.RepoCleanBg,
		})
	}

}
