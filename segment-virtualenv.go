package main

import (
	"os"
	"path"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentVirtualEnv(p *powerline) {
	var env string
	if env == "" {
		env, _ = os.LookupEnv("VIRTUAL_ENV")
	}
	if env == "" {
		env, _ = os.LookupEnv("CONDA_ENV_PATH")
	}
	if env == "" {
		env, _ = os.LookupEnv("CONDA_DEFAULT_ENV")
	}
	if env == "" {
		return
	} else {
		envName := path.Base(env)
		p.appendSegment("venv", pwl.Segment{
			Content:    envName,
			Foreground: p.theme.VirtualEnvFg,
			Background: p.theme.VirtualEnvBg,
		})
	}
}
