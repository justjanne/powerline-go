package main

import (
	"log"
	"os/exec"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentGCP(p *powerline) []pwl.Segment {
	out, err := exec.Command("gcloud", "config", "list", "project", "--format", "value(core.project)").Output()
	if err != nil {
		log.Fatal(err)
	}

	project := strings.TrimSuffix(string(out), "\n")
	if project == "" {
		return []pwl.Segment{}
	}
	return []pwl.Segment{{
		Name:       "gcp",
		Content:    project,
		Foreground: p.theme.GCPFg,
		Background: p.theme.GCPBg,
	}}
}
