package main

import (
	"os"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentAWS(p *powerline) []pwl.Segment {
	profile := os.Getenv("AWS_PROFILE")
	region := os.Getenv("AWS_DEFAULT_REGION")
	var content string = ""

	if len(profile) == 0 {
		profile = os.Getenv("AWS_VAULT")
	}

	if len(region) == 0 {
		region = os.Getenv("AWS_REGION")
	}

	if len(region) > 0 {
		content = "(" + region + ")"
	}

	if len(profile) > 0 {
		content = profile +" "+ content
	}

	if len(content) == 0 {
		return []pwl.Segment{}
	}

	return []pwl.Segment{{
		Name:       "aws",
		Content:    content,
		Foreground: p.theme.AWSFg,
		Background: p.theme.AWSBg,
	}}
}
