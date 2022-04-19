package main

import (
	"encoding/json"
	pwl "github.com/justjanne/powerline-go/powerline"
	"io/ioutil"
	"os"
	"path/filepath"
)

type DockerContextConfig struct {
	CurrentContext string `json:"currentContext"`
}

func segmentDockerContext(p *powerline) []pwl.Segment {
	context := "default"
	home, _ := os.LookupEnv("HOME")
	contextFolder := filepath.Join(home, ".docker", "contexts")
	configFile := filepath.Join(home, ".docker", "config.json")
	contextEnvVar := os.Getenv("DOCKER_CONTEXT")

	if contextEnvVar != "" {
		context = contextEnvVar
	} else {
		stat, err := os.Stat(contextFolder)
		if err == nil && stat.IsDir() {
			dockerConfigFile, err := ioutil.ReadFile(configFile)
			if err == nil {
				var dockerConfig DockerContextConfig
				err = json.Unmarshal(dockerConfigFile, &dockerConfig)
				if err == nil && dockerConfig.CurrentContext != "" {
					context = dockerConfig.CurrentContext
				}
			}
		}
	}

	// Don’t show the default context
	if context == "default" {
		return []pwl.Segment{}
	}

	return []pwl.Segment{{
		Name:       "docker-context",
		Content:    "🐳" + context,
		Foreground: p.theme.PlEnvFg,
		Background: p.theme.PlEnvBg,
	}}
}
