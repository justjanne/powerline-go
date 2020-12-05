package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
)

func segmentGoenv(p *powerline) []pwl.Segment {
	env, _ := os.LookupEnv("GOENV_VERSION")
	if env != "" {
		return []pwl.Segment{{
			Name:       "goenv",
			Content:    env,
			Foreground: p.theme.GoEnvFg,
			Background: p.theme.GoEnvBg,
		}}
	}
	return []pwl.Segment{}
}
