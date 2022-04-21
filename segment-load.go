package main

import (
	"fmt"
	pwl "github.com/justjanne/powerline-go/powerline"
	"runtime"

	"github.com/shirou/gopsutil/v3/load"
)

func segmentLoad(p *powerline) []pwl.Segment {
	c := runtime.NumCPU()
	a, err := load.Avg()
	if err != nil {
		return []pwl.Segment{}
	}
	bg := p.theme.LoadBg

	load := a.Load5
	switch p.theme.LoadAvgValue {
	case 1:
		load = a.Load1
	case 15:
		load = a.Load15
	}

	if load > float64(c)*p.theme.LoadThresholdBad {
		bg = p.theme.LoadHighBg
	}

	return []pwl.Segment{{
		Name:       "load",
		Content:    fmt.Sprintf("%.2f", a.Load5),
		Foreground: p.theme.LoadFg,
		Background: bg,
	}}
}
