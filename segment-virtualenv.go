package main

import (
	"os"
	"fmt"
	"path"
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
		p.appendSegment(segment{
			content:    fmt.Sprintf(" %s ", envName),
			foreground: p.theme.VirtualEnvFg,
			background: p.theme.VirtualEnvBg,
		})
	}
}
