package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"os/exec"
	"strings"
)

func runRbenvCommand(cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)
	out, err := command.Output()
	return string(out), err
}

func segmentRbenv(p *powerline) []pwl.Segment {
	out, err := runRbenvCommand("rbenv", "version")
	
	if err == nil {
		items := strings.Split(out, " ")
		if len(items) > 1 {
			return []pwl.Segment{{
				Name:       "rbenv",
				Content:    items[0],
				Foreground: p.theme.TimeFg,
				Background: p.theme.TimeBg,
			}}	
		}
	}

	return []pwl.Segment{}
}
