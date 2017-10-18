package main

import (
	"crypto/md5"
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
		hostPrompt = hostName

		hash := getMd5(hostName)
		background = hash[0]
		foreground = p.theme.HostnameColorizedFgMap[background]
	} else {
		if *p.args.Shell == "bash" {
			hostPrompt = "\\h"
		} else if *p.args.Shell == "zsh" {
			hostPrompt = "%m"
		} else {
			hostPrompt = getHostName()
		}

		foreground = p.theme.HostnameFg
		background = p.theme.HostnameBg
	}

	p.appendSegment("host", segment{
		content:    hostPrompt,
		foreground: foreground,
		background: background,
	})
}
