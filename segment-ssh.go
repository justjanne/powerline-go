package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
)

func segmentSSH(p *powerline) []pwl.Segment {
	sshClient, _ := os.LookupEnv("SSH_CLIENT")
	if sshClient != "" {
		return []pwl.Segment{{
			Name:       "ssh",
			Content:    p.symbolTemplates.Network,
			Foreground: p.theme.SSHFg,
			Background: p.theme.SSHBg,
		}}
	}
	return []pwl.Segment{}
}
