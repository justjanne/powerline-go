package main

import (
	"fmt"
	"time"
)

func segmentTime(p *powerline) {
	p.appendSegment("time", segment{
		content:    fmt.Sprintf(" %s ", time.Now().Format("15:04:05")),
		foreground: p.theme.TimeFg,
		background: p.theme.TimeBg,
	})
}
