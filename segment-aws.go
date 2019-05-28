package main

import (
	"os"
)

func segmentAWS(p *powerline) {
	profile := os.Getenv("AWS_PROFILE")
	region := os.Getenv("AWS_DEFAULT_REGION")
	if profile != "" {
		var r string
		if region != "" {
			r = " (" + region + ")"
		}
		p.appendSegment("aws", segment{
			content:    profile + r,
			foreground: p.theme.AWSFg,
			background: p.theme.AWSBg,
		})
	}
}
