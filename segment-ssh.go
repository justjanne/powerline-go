package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
)

func segmentSSH(p *powerline) {
	sshClient, _ := os.LookupEnv("SSH_CLIENT")
	if sshClient != "" {
		p.appendSegment("ssh", pwl.Segment{
			Content:    p.symbolTemplates.Network,
			Foreground: p.theme.SSHFg,
			Background: p.theme.SSHBg,
		})
	}
}
