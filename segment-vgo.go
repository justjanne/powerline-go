package main

import (
	"os"
)

func segmentVirtualGo(p *powerline) {
	var env string
	if env == "" {
		env, _ = os.LookupEnv("VIRTUALGO")
	}
	if env == "" {
		return
	}
	p.appendSegment("vgo", segment{
		content:    env,
		foreground: p.theme.VirtualGoFg,
		background: p.theme.VirtualGoBg,
	})
}
