package main

import (
	"os"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentDotEnv(p *powerline) []pwl.Segment {
	files := []string{".env", ".envrc"}
	dotEnv := false
	for _, file := range files {
		stat, err := os.Stat(file)
		if err == nil && !stat.IsDir() {
			dotEnv = true
			break
		}
	}
	if !dotEnv {
		return []pwl.Segment{}
	}
	return []pwl.Segment{{
		Name:       "dotenv",
		Content:    "\u2235",
		Foreground: p.theme.DotEnvFg,
		Background: p.theme.DotEnvBg,
	}}
}
