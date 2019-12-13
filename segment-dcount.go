package main

import (
	"fmt"
	"io/ioutil"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentDcount(p *powerline) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return
	}
	nFiles := len(files)

	if nFiles > 0 {
		p.appendSegment("dcount", pwl.Segment{
			Content:    fmt.Sprintf("%d", nFiles),
			Foreground: p.theme.DirCountFg,
			Background: p.theme.DirCountBg,
		})
	}
}
