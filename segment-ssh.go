package main

import (
	"os"
	"fmt"
)

func segmentSsh(p *powerline) {
	sshClient, _ := os.LookupEnv("SSH_CLIENT")
	if sshClient != "" {
		p.appendSegment(segment{
			content:    fmt.Sprintf(" %s ", p.symbolTemplates.Network),
			foreground: p.theme.SshFg,
			background: p.theme.SshBg,
		})
	}
}
