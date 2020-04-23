package main

import (
	"log"
	"os/exec"
	"strings"
)

func segmentGCP(p *powerline) {
	out, err := exec.Command("gcloud", "config", "list", "project", "--format", "value(core.project)").Output()
	if err != nil {
		log.Fatal(err)
	}

	project := strings.TrimSuffix(string(out), "\n")
	if project != "" {
		p.appendSegment("gcp", segment{
			content:    project,
			foreground: p.theme.GCPFg,
			background: p.theme.GCPBg,
		})
	}
}
