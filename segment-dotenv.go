package main

import (
	"os"
)

func segmentDotEnv(p *powerline) {
	stat, err := os.Stat(".env")
	if err == nil && !stat.IsDir() {
		p.appendSegment("dotenv", segment{
			content:    " \u2235 ",
			foreground: p.theme.DotEnvFg,
			background: p.theme.DotEnvBg,
		})
	}
}
