package main

import (
	"net/url"
	"os"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentDocker(p *powerline) []pwl.Segment {
	var docker string
	dockerMachineName, _ := os.LookupEnv("DOCKER_MACHINE_NAME")
	dockerHost, _ := os.LookupEnv("DOCKER_HOST")

	if dockerMachineName != "" {
		docker = dockerMachineName
	} else if dockerHost != " " {
		u, err := url.Parse(dockerHost)
		if err == nil {
			docker = u.Host
		}
	}

	if docker != "" {
		return []pwl.Segment{{
			Name:       "docker",
			Content:    docker,
			Foreground: p.theme.DockerMachineFg,
			Background: p.theme.DockerMachineBg,
		}}
	}
	return []pwl.Segment{}
}
