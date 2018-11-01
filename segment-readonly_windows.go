// +build windows

package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
)

func segmentPerms(p *powerline) {
	cwd := p.cwd
	if cwd == "" {
		cwd, _ = os.LookupEnv("PWD")
	}

	const W_USR = 0002
	// Check user's permissions on directory in a portable but probably slower way
	fileInfo, _ := os.Stat(cwd)
	if fileInfo.Mode()&W_USR != W_USR {
		p.appendSegment("perms", pwl.Segment{
			Content:    p.symbolTemplates.Lock,
			Foreground: p.theme.ReadonlyFg,
			Background: p.theme.ReadonlyBg,
		})
	}
}
