package main

import (
	"os"
	"path"
)

func segmentPerlbrew(p *powerline) {
	env, _ := os.LookupEnv("PERLBREW_PERL")
	if env == "" {
		return
	}

	envName := path.Base(env)
	p.appendSegment("perlbrew", segment{
		content:    envName,
		foreground: p.theme.PerlbrewFg,
		background: p.theme.PerlbrewBg,
	})
}
