package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"log"
	"os/exec"
	"strings"
)

func segmentGCP(p *powerline) []pwl.Segment {
	out, err := exec.Command("gcloud", "config", "list", "project", "--format", "value(core.project)").Output()
	if err != nil {
		log.Fatal(err)
	}

	project := strings.TrimSuffix(string(out), "\n")
	if project != "" {
		return []pwl.Segment{{
			Name:       "gcp",
			Content:    project,
			Foreground: p.theme.GCPFg,
			Background: p.theme.GCPBg,
		}}
	}

	return []pwl.Segment{}
}
