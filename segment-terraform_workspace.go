package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"io/ioutil"
	"os"
)

const wsFile = "./.terraform/environment"

func segmentTerraformWorkspace(p *powerline) []pwl.Segment {
	stat, err := os.Stat(wsFile)

	if err == nil && !stat.IsDir() {
		workspace, err := ioutil.ReadFile(wsFile)
		if err == nil {
			return []pwl.Segment{{
				Name:       "terraform-workspace",
				Content:    string(workspace),
				Foreground: p.theme.TFWsFg,
				Background: p.theme.TFWsBg,
			}}

		}
	}
	return []pwl.Segment{}
}
