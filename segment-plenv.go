package main

import (
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
		p.appendSegment("plenv", segment{
			content:    envName,
			foreground: p.theme.PlEnvFg,
			background: p.theme.PlEnvBg,
		})
	}
}
