package main

import (
	"os"
)

var editSegment = segment{
	content:    "\u270E",
	foreground: 15,
	background: 161,
}

var commandSegment = segment{
	content:    "\u26CF",
	foreground: 15,
	background: 31,
}

func segmentKeymap(p *powerline) {
	var keymap string
	if keymap == "" {
		keymap, _ = os.LookupEnv("KEYMAP_POWERLINE")
	}
	if keymap == "main" {
		p.appendSegment("keymap", editSegment)
	}
	if keymap == "vicmd" {
		p.appendSegment("keymap", commandSegment)
	}

}
