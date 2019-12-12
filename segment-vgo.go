package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
)

func segmentVirtualGo(p *powerline) []pwl.Segment {
	var env string
	if env == "" {
		env, _ = os.LookupEnv("VIRTUALGO")
	}
	segments := []pwl.Segment{}
	if env == "" {
		return segments
	} else {
		segments = append(segments, pwl.Segment{
			Name:       "vgo",
			Content:    env,
			Foreground: p.theme.VirtualGoFg,
			Background: p.theme.VirtualGoBg,
		})
	}
	segments = append(segments, pwl.Segment{
		Name:       "vgo",
		Content:    env,
		Foreground: p.theme.VirtualGoFg,
		Background: p.theme.VirtualGoBg,
	})
	return segments
}
