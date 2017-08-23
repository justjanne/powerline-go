package main

import (
	"crypto/md5"
	"fmt"
	"os"
	"strings"
)

func getHostName() string {
	fullyQualifiedDomainName, _ := os.Hostname()
	return strings.SplitN(fullyQualifiedDomainName, ".", 2)[0]
}

func getMd5(text string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hasher.Sum(nil)
}

func segmentHost(p *powerline) {
	var hostPrompt string
	var foreground, background uint8
	if *p.args.ColorizeHostname {
		hostName := getHostName()
		hostPrompt = fmt.Sprintf(" %s ", hostName)

		hash := getMd5(hostName)
		background = hash[0]
		foreground = p.theme.HostnameColorizedFgMap[background]
	} else {
		if *p.args.Shell == "bash" {
			hostPrompt = " \\h "
		} else if *p.args.Shell == "zsh" {
			hostPrompt = " %m "
		} else {
			hostPrompt = fmt.Sprintf(" %s ", getHostName())
		}

		foreground = p.theme.HostnameFg
		background = p.theme.HostnameBg
	}

	p.appendSegment(segment{
		content:    hostPrompt,
		foreground: foreground,
		background: background,
	})
}
