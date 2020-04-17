package main

import (
	"crypto/md5"
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
	"strconv"
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

		foregroundEnvStr := os.Getenv("PLGO_HOSTNAMEFG")
		backgroundEnvStr := os.Getenv("PLGO_HOSTNAMEBG")
		foregroundEnv, foregroundErr := strconv.ParseInt(foregroundEnvStr, 0, 64)
		backgroundEnv, backgroundErr := strconv.ParseInt(backgroundEnvStr, 0, 64)

		if foregroundErr == nil || backgroundErr == nil {
			if foregroundErr != nil {
				foreground = p.theme.HostnameFg
			} else {
				foreground = uint8(foregroundEnv)
			}
			if backgroundErr != nil {
				background = p.theme.HostnameBg
			} else {
				background = uint8(backgroundEnv)
			}
		} else {
			hash := getMd5(hostName)
			background = hash[0]

			if foregroundMap, exists := p.theme.HostnameColorizedFgMap[background]; !exists {
				foreground = p.theme.HostnameFg
				background = p.theme.HostnameBg
			} else {
				foreground = foregroundMap
			}
		}
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
