package main

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"gopkg.in/ini.v1"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentVirtualEnv(p *powerline) []pwl.Segment {
	var env string
	if env == "" {
		env, _ = os.LookupEnv("VIRTUAL_ENV")
		if env != "" {
			cfg, err := ini.Load(path.Join(env, "pyvenv.cfg"))
			if err == nil {
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
	pyenv := false
	if env == "" {
		env, _ = os.LookupEnv("PYENV_VERSION")
		pyenv = true
	}
	if env == "" && os.Getenv("PYENV_ROOT") != "" {
		if out, err := exec.Command("pyenv", "version-name").Output(); err == nil {
			env = strings.SplitN(strings.TrimSpace(string(out)), ":", 2)[0]
			pyenv = true
		}
	}
	if env == "" {
		return []pwl.Segment{}
	}
	if pyenv {
		if out, err := exec.Command("pyenv", "global").Output(); err == nil {
			if env == strings.SplitN(strings.TrimSpace(string(out)), ":", 2)[0] {
				return []pwl.Segment{}
			}
		}
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
