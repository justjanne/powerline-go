package main

import (
	"net/url"
	"os"
)

func segmentDocker(p *powerline) {
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
		p.appendSegment("docker", segment{
			content:    docker,
			foreground: p.theme.DockerMachineFg,
			background: p.theme.DockerMachineBg,
		})
	}
}
