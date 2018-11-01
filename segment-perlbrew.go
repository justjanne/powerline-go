package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
	"path"
)

func segmentPerlbrew(p *powerline) {
	env, _ := os.LookupEnv("PERLBREW_PERL")
	if env == "" {
		return
	}

	envName := path.Base(env)
	p.appendSegment("perlbrew", pwl.Segment{
		Content:    envName,
		Foreground: p.theme.PerlbrewFg,
		Background: p.theme.PerlbrewBg,
	})
}
