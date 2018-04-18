package main

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/load"
)

func segmentLoad(p *powerline) {
	c := runtime.NumCPU()
	a, err := load.Avg()
	if err != nil {
		return
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

	p.appendSegment("load", segment{
		content:    fmt.Sprintf("%.2f", a.Load5),
		foreground: p.theme.LoadFg,
		background: bg,
	})
}
