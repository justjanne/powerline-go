package main

import (
	"fmt"
	"os"
)

func segmentSsh(p *powerline) {
	sshClient, _ := os.LookupEnv("SSH_CLIENT")
	if sshClient != "" {
		p.appendSegment("ssh", segment{
			content:    fmt.Sprintf(" %s ", p.symbolTemplates.Network),
			foreground: p.theme.SshFg,
			background: p.theme.SshBg,
		})
	}
}
