package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"strconv"
)

func segmentJobs(p *powerline) []pwl.Segment {
	if *p.args.Jobs > 0 {
		return []pwl.Segment{{
			Name:       "jobs",
			Content:    strconv.Itoa(*p.args.Jobs),
			Foreground: p.theme.JobsFg,
			Background: p.theme.JobsBg,
		}}
	}
	return []pwl.Segment{}
}
