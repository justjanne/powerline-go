package main

import (
	"os"
	"fmt"
	"strings"
)

func segmentHost(p *powerline) {
	var hostPrompt string
	if *p.args.Shell == "bash" {
		hostPrompt = " \\h "
	} else if *p.args.Shell == "zsh" {
		hostPrompt = " %m "
	} else {
		fullyQualifiedDomainName, _ := os.Hostname()
		hostname := strings.SplitN(fullyQualifiedDomainName, ".", 1)[0]
		hostPrompt = fmt.Sprintf(" %s ", hostname)
	}

	p.appendSegment(segment{
		content:    hostPrompt,
		foreground: p.theme.HostnameFg,
		background: p.theme.HostnameBg,
	})
}
