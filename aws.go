package main

import (
	"os"
)

func segmentAWS(p *powerline) {
	profile, _ := os.LookupEnv("AWS_PROFILE")
	if profile != "" {
		p.appendSegment("aws", segment{
			content:    profile,
			foreground: p.theme.AWSFg,
			background: p.theme.AWSBg,
		})
	}
}
