package main

import (
	"encoding/json"
	"os/exec"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentPlugin(p *powerline, plugin string) bool {
	output, err := exec.Command("powerline-go-" + plugin).Output()
	if err != nil {
		return false
	}
	segments := []pwl.Segment{}
	err = json.Unmarshal(output, &segments)
	if err != nil {
		// The plugin was found but no valid data was returned. Ignore it
		return true
	}
	for _, s := range segments {
		p.appendSegment(plugin, s)
	}
	return true
}
