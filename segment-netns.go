package main

import (
	"os/exec"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentNetns(p *powerline) []pwl.Segment {
	output, err := exec.Command("ip", "netns", "identify").Output()
	if err != nil || len(output) == 0 {
		return []pwl.Segment{}
	}
	outs := strings.TrimSpace(string(output)) // may consist of just a linefeed
	if len(outs) == 0 {
		return []pwl.Segment{}
	}
	return []pwl.Segment{{
		Name:       "netns",
		Content:    outs,
		Foreground: p.theme.NetnsFg,
		Background: p.theme.NetnsBg,
	}}
}
