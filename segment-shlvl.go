package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
	"strconv"
)

func segmentShlvl(p *powerline) []pwl.Segment {

	level, _ := os.LookupEnv("SHLVL")
	leveli, err := strconv.Atoi(level)

	if err != nil || leveli < 1 {
		return []pwl.Segment{}
	}
	return []pwl.Segment{{
		Name:       "shlvl",
		Content:    level,
		Foreground: p.theme.ShLvlFg,
		Background: p.theme.ShLvlBg,
	}}
}
