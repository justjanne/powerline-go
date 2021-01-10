package main

import (
	"os"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentPlEnv(p *powerline) []pwl.Segment {
	env, _ := os.LookupEnv("PLENV_VERSION")
	if env == "" {
		return []pwl.Segment{}
	}
	return []pwl.Segment{{
		Name:       "plenv",
		Content:    env,
		Foreground: p.theme.PlEnvFg,
		Background: p.theme.PlEnvBg,
	}}
}
