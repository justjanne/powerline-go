package main

import (
	"os"
)

func segmentDotEnv(p *powerline) {
	files := []string{".env", ".envrc"}
	dotEnv := false
	for _, file := range files {
		stat, err := os.Stat(file)
		if err == nil && !stat.IsDir() {
			dotEnv = true
			break
		}
	}
	if dotEnv {
		p.appendSegment("dotenv", segment{
			content:    "\u2235",
			foreground: p.theme.DotEnvFg,
			background: p.theme.DotEnvBg,
		})
	}
}
