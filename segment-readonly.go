// +build !windows

package main

import (
	"os"

	"golang.org/x/sys/unix"
)

func segmentPerms(p *powerline) {
	cwd := p.cwd
	if cwd == "" {
		cwd, _ = os.LookupEnv("PWD")
	}
	if unix.Access(cwd, unix.W_OK) != nil {
		p.appendSegment("perms", segment{
			content:    p.symbolTemplates.Lock,
			foreground: p.theme.ReadonlyFg,
			background: p.theme.ReadonlyBg,
		})
	}
}
