package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
	"path"
)

func segmentShEnv(p *powerline) []pwl.Segment {
	var env string
	if env == "" {
		env, _ = os.LookupEnv("SHENV_VERSION")
	}
	if env != "" {
		envName := path.Base(env)
		return []pwl.Segment{{
			Name: "shenv",
			Content:    envName,
			Foreground: p.theme.ShEnvFg,
			Background: p.theme.ShEnvBg,
		}}
	}

	return []pwl.Segment{}
}
