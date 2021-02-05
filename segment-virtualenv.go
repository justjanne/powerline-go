package main

import (
	"bufio"
	"os"
	"path"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentVirtualEnv(p *powerline) []pwl.Segment {
	var env string
	if env == "" {
		env, _ = os.LookupEnv("VIRTUAL_ENV")
		if env != "" {
			file, err := os.Open(path.Join(env, "pyvenv.cfg"))
			if err == nil {
				defer file.Close()
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					if strings.Contains(scanner.Text(), "prompt") {
						var prompt = strings.Split(scanner.Text(), "=")[1]
						prompt = strings.TrimSpace(prompt)
						prompt = strings.Trim(prompt, "'")
						env = prompt
						break
					}
				}
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
		Content:    envName,
		Foreground: p.theme.VirtualEnvFg,
		Background: p.theme.VirtualEnvBg,
	}}
}
