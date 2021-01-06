package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
)

func segmentShEnv(p *powerline) []pwl.Segment {
	env, _ := os.LookupEnv("SHENV_VERSION")
	if env != "" {
		return []pwl.Segment{{
			Name:       "shenv",
			Content:    env,
			Foreground: p.theme.ShEnvFg,
			Background: p.theme.ShEnvBg,
		}}
	}

	return []pwl.Segment{}
}
