package main

import (
	"net/url"
	"os"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentWSL(p *powerline) []pwl.Segment {
	var WSL string
	WSLMachineName, _ := os.LookupEnv("WSL_DISTRO_NAME")
	WSLHost, _ := os.LookupEnv("NAME")

	if WSLMachineName != "" {
		WSL = WSLMachineName
	} else if WSLHost != " " {
		u, err := url.Parse(WSLHost)
		if err == nil {
			WSL = u.Host
		}
	}

	if WSL != "" {
		return []pwl.Segment{{
			Name:       "WSL",
			Content:    WSL,
			Foreground: p.theme.WSLMachineFg,
			Background: p.theme.WSLMachineBg,
		}}
	}
	return []pwl.Segment{}
}
