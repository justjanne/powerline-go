package main

import (
	"os"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentMise(p *powerline) []pwl.Segment {
	content := os.Getenv("MISE_ENV")
	if content == "" {
		return []pwl.Segment{}
	}
	return []pwl.Segment{{
		Name:       "mise",
		Content:    p.symbols.MiseIndicator + " " + content,
		Foreground: p.theme.MiseFg,
		Background: p.theme.MiseBg,
	}}
}
