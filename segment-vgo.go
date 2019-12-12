package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
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
	p.appendSegment("vgo", pwl.Segment{
		Content:    env,
		Foreground: p.theme.VirtualGoFg,
		Background: p.theme.VirtualGoBg,
	})
}
