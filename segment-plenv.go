package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
	"path"
)

func segmentPlEnv(p *powerline) []pwl.Segment {
	var env string
	if env == "" {
		env, _ = os.LookupEnv("PLENV_VERSION")
	}
	if env != "" {
		envName := path.Base(env)
		return []pwl.Segment{{
			Name:       "plenv",
			Content:    envName,
			Foreground: p.theme.PlEnvFg,
			Background: p.theme.PlEnvBg,
		}}
	}
	return []pwl.Segment{}
}
