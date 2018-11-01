package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
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
		p.appendSegment("aws", pwl.Segment{
			Content:    profile + r,
			Foreground: p.theme.AWSFg,
			Background: p.theme.AWSBg,
		})
	}
}
