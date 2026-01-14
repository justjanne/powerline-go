package main

import (
	"os"
	"os/exec"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
)

const tfDirName = "./.terraform"

func runTfCommand(cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)
	out, err := command.Output()
	return string(out), err
}

func segmentTerraformWorkspace(p *powerline) []pwl.Segment {
	if _, err := os.Stat(tfDirName); os.IsNotExist(err) {
		return []pwl.Segment{}
	}

	workspace, err := runTfCommand("terraform", "workspace", "show")
	if err != nil {
		return []pwl.Segment{}
	}

	return []pwl.Segment{{
		Name:       "terraform-workspace",
		Content:    strings.TrimSpace(workspace),
		Foreground: p.theme.TFWsFg,
		Background: p.theme.TFWsBg,
	}}
}
