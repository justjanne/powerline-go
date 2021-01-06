package main

import (
	"os"
	"path"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentVirtualEnv(p *powerline) []pwl.Segment {
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
		env, _ = os.LookupEnv("PYENV_VERSION")
	}
	segments := []pwl.Segment{}
	if env != "" {
		envName := path.Base(env)
		if p.cfg.VenvNameSizeLimit > 0 && len(envName) > p.cfg.VenvNameSizeLimit {
			envName = p.symbols.VenvIndicator
		}
		segments = append(segments, pwl.Segment{
			Name:       "venv",
			Content:    envName,
			Foreground: p.theme.VirtualEnvFg,
			Background: p.theme.VirtualEnvBg,
		})
	}
	return segments
}
