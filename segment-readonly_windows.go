// +build windows

package main

import (
	"fmt"
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
		p.appendSegment("perms", segment{
			content:    fmt.Sprintf(" %s ", p.symbolTemplates.Lock),
			foreground: p.theme.ReadonlyFg,
			background: p.theme.ReadonlyBg,
		})
	}
}
