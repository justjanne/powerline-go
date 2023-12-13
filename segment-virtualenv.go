package main

import (
	"os"
	"path"

	"gopkg.in/ini.v1"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentVirtualEnv(p *powerline) []pwl.Segment {
	var env string
	if env == "" {
		// honor the $VIRTUAL_ENV_PATH first
		env, _ = os.LookupEnv("VIRTUAL_ENV_PATH")
	}
	if env == "" {
		env, _ = os.LookupEnv("VIRTUAL_ENV")
		if env != "" {
			cfg, err := ini.Load(path.Join(env, "pyvenv.cfg"))
			// in the case of a "prompt" value not being set in cfg,
			// Key() will create an empty value and return it. this
			// obliterates the env derived from VIRTUAL_ENV
			if err == nil && cfg.Section("").HasKey("prompt") {
				env = cfg.Section("").Key("prompt").String()
			}
		}
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
	if env == "" {
		return []pwl.Segment{}
	}
	envName := path.Base(env)
	if p.cfg.VenvNameSizeLimit > 0 && len(envName) > p.cfg.VenvNameSizeLimit {
		envName = p.symbols.VenvIndicator
	}

	return []pwl.Segment{{
		Name:       "venv",
		Content:    escapeVariables(p, envName),
		Foreground: p.theme.VirtualEnvFg,
		Background: p.theme.VirtualEnvBg,
	}}
}
