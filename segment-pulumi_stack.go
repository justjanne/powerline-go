package main

import (
	"os"
	"os/exec"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
)

const confFile = "./Pulumi.yaml"

func segmentPulumiStack(p *powerline) []pwl.Segment {
	stat, err := os.Stat(confFile)
	if err != nil {
		return []pwl.Segment{}
	}
	if stat.IsDir() {
		return []pwl.Segment{}
	}
	command := exec.Command("pulumi", "stack", "--show-name")
	stack, err := command.Output()
	if err != nil {
		return []pwl.Segment{}
	}
	return []pwl.Segment{{
		Name:       "pulumi-stack",
		Content:    strings.TrimSuffix(string(stack), "\n"),
		Foreground: p.theme.PulumiStackFg,
		Background: p.theme.PulumiStackBg,
	}}
}
