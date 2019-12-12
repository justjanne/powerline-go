package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
	"path"
)

func segmentShEnv(p *powerline) {
	var env string
	if env == "" {
		env, _ = os.LookupEnv("SHENV_VERSION")
	}
	if env == "" {
		return
	}
	envName := path.Base(env)
	p.appendSegment("shenv", pwl.Segment{
		Content:    envName,
		Foreground: p.theme.ShEnvFg,
		Background: p.theme.ShEnvBg,
	})
}
