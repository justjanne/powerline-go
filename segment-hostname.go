package main

import (
	"crypto/md5"
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
	"strconv"
	"strings"
)

func getHostName(fullyQualifiedDomainName string, keepFqdnHostname bool) string {
	if keepFqdnHostname {
		return fullyQualifiedDomainName
	}
	return strings.SplitN(fullyQualifiedDomainName, ".", 2)[0]
}

func getMd5(text string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hasher.Sum(nil)
}

func segmentHost(p *powerline) []pwl.Segment {
	var hostPrompt string
	var foreground, background uint8

	if p.cfg.HostnameOnlyIfSSH {
		if os.Getenv("SSH_CLIENT") == "" {
			// It's not an ssh connection do nothing
			return []pwl.Segment{}
		}
	}

	if p.cfg.ColorizeHostname {
		hostName := getHostName(p.hostname, p.cfg.FqdnHostname)
		hostPrompt = hostName

		foregroundEnv, foregroundEnvErr := strconv.ParseUint(os.Getenv("PLGO_HOSTNAMEFG"), 0, 8)
		backgroundEnv, backgroundEnvErr := strconv.ParseUint(os.Getenv("PLGO_HOSTNAMEBG"), 0, 8)
		if foregroundEnvErr == nil && backgroundEnvErr == nil {
			foreground = uint8(foregroundEnv)
			background = uint8(backgroundEnv)
		} else {
			hash := getMd5(hostName)
			background = hash[0] % 128
			foreground = p.theme.HostnameColorizedFgMap[background]
		}
	} else {
		if p.cfg.Shell == "bash" {
			hostPrompt = "\\h"
		} else if p.cfg.Shell == "zsh" {
			hostPrompt = "%m"
		} else {
			hostPrompt = getHostName(p.hostname, p.cfg.FqdnHostname)
		}

		foreground = p.theme.HostnameFg
		background = p.theme.HostnameBg
	}

	return []pwl.Segment{{
		Name:       "host",
		Content:    hostPrompt,
		Foreground: foreground,
		Background: background,
	}}
}
