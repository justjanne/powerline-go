// +build windows

package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
)

func segmentPerms(p *powerline) []pwl.Segment {
	cwd := p.cwd
	const W_USR = 0002
	// Check user's permissions on directory in a portable but probably slower way
	fileInfo, _ := os.Stat(cwd)
	if fileInfo.Mode()&W_USR != W_USR {
		return []pwl.Segment{{
			Name:       "perms",
			Content:    p.symbols.Lock,
			Foreground: p.theme.ReadonlyFg,
			Background: p.theme.ReadonlyBg,
		}}
	}
	return []pwl.Segment{}
}
