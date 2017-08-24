package main

import (
	"fmt"
	"os"
)

func segmentDocker(p *powerline) {
	dockerMachineName, _ := os.LookupEnv("DOCKER_MACHINE_NAME")
	if dockerMachineName != "" {
		p.appendSegment("docker", segment{
			content:    fmt.Sprintf(" %s ", dockerMachineName),
			foreground: p.theme.DockerMachineFg,
			background: p.theme.DockerMachineBg,
		})
	}
}
