package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
)

func segmentPlEnv(p *powerline) []pwl.Segment {
	env, _ := os.LookupEnv("PLENV_VERSION")
	if env != "" {
		return []pwl.Segment{{
			Name:       "plenv",
			Content:    env,
			Foreground: p.theme.PlEnvFg,
			Background: p.theme.PlEnvBg,
		}}
	}
	return []pwl.Segment{}
}
