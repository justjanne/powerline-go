package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
	"path"
)

func segmentPlEnv(p *powerline) {
	var env string
	if env == "" {
		env, _ = os.LookupEnv("PLENV_VERSION")
	}
	if env == "" {
		return
	} else {
		envName := path.Base(env)
		p.appendSegment("plenv", pwl.Segment{
			Content:    envName,
			Foreground: p.theme.PlEnvFg,
			Background: p.theme.PlEnvBg,
		})
	}
}
