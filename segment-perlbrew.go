package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
	"path"
)

func segmentPerlbrew(p *powerline) []pwl.Segment {
	env, _ := os.LookupEnv("PERLBREW_PERL")
	if env == "" {
		return []pwl.Segment{}
	}

	envName := path.Base(env)
	return []pwl.Segment{{
		Name:       "perlbrew",
		Content:    envName,
		Foreground: p.theme.PerlbrewFg,
		Background: p.theme.PerlbrewBg,
	}}
}
