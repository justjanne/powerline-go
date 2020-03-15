package main

import (
	"crypto/md5"
	pwl "github.com/justjanne/powerline-go/powerline"
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

	if *p.args.HostnameOnlyIfSSH {
		if os.Getenv("SSH_CLIENT") == "" {
			// It's not an ssh connection do nothing
			return
		}
	}

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

	p.appendSegment("host", pwl.Segment{
		Content:    hostPrompt,
		Foreground: foreground,
		Background: background,
	})
}
