package main

import (
	"fmt"
	"os"
)

func segmentAWS(p *powerline) {
	profile, _ := os.LookupEnv("AWS_PROFILE")
	if profile != "" {
		p.appendSegment("aws", segment{
			content:    fmt.Sprintf(" %s ", profile),
			foreground: p.theme.AWSFg,
			background: p.theme.AWSBg,
		})
	}
}
