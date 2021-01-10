package main

import (
	"os"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentSSH(p *powerline) []pwl.Segment {
	sshClient, _ := os.LookupEnv("SSH_CLIENT")
	if sshClient == "" {
		return []pwl.Segment{}
	}
	var networkIcon string
	if p.cfg.SshAlternateIcon {
		networkIcon = p.symbols.NetworkAlternate
	} else {
		networkIcon = p.symbols.Network
	}

	return []pwl.Segment{{
		Name:       "ssh",
		Content:    networkIcon,
		Foreground: p.theme.SSHFg,
		Background: p.theme.SSHBg,
	}}
}
