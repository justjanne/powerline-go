package main

import (
	"os"
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
	}
	envName := path.Base(env)
	p.appendSegment("venv", segment{
		content:    envName,
		foreground: p.theme.VirtualEnvFg,
		background: p.theme.VirtualEnvBg,
	})
}
