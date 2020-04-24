package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"time"
)

func segmentTime(p *powerline) {
	p.appendSegment("time", pwl.Segment{
		Content:    time.Now().Format("15:04:05"),
		Foreground: p.theme.TimeFg,
		Background: p.theme.TimeBg,
	})
}
