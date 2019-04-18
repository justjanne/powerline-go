package main

import (
	"os"
)

func segmentSSH(p *powerline) {
	sshClient, _ := os.LookupEnv("SSH_CLIENT")
	if sshClient != "" {
		p.appendSegment("ssh", segment{
			content:    p.symbolTemplates.Network,
			foreground: p.theme.SSHFg,
			background: p.theme.SSHBg,
		})
	}
}
