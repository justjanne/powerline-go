// +build !windows

package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"

	"golang.org/x/sys/unix"
)

func segmentPerms(p *powerline) {
	cwd := p.cwd
	if cwd == "" {
		cwd, _ = os.LookupEnv("PWD")
	}
	if unix.Access(cwd, unix.W_OK) != nil {
		p.appendSegment("perms", pwl.Segment{
			Content:    p.symbolTemplates.Lock,
			Foreground: p.theme.ReadonlyFg,
			Background: p.theme.ReadonlyBg,
		})
	}
}
