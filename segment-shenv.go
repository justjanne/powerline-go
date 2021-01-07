package main

import (
	"os"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentShEnv(p *powerline) []pwl.Segment {
	env, _ := os.LookupEnv("SHENV_VERSION")
	if env == "" {
		return []pwl.Segment{}
	}
	return []pwl.Segment{{
		Name:       "shenv",
		Content:    env,
		Foreground: p.theme.ShEnvFg,
		Background: p.theme.ShEnvBg,
	}}
}
