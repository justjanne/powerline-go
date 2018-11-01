package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	pwl "github.com/justjanne/powerline-go/powerline"
)

const pkgfile = "./package.json"

type packageJSON struct {
	Version string `json:"version"`
}

func segmentNode(p *powerline) {
	stat, err := os.Stat(pkgfile)
	if err == nil && !stat.IsDir() {
		pkg := packageJSON{"!"}
		raw, err := ioutil.ReadFile(pkgfile)
		if err == nil {
			err = json.Unmarshal(raw, &pkg)
			if err == nil {
				p.appendSegment("node-version", pwl.Segment{
					Content:    pkg.Version + " \u2B22",
					Foreground: p.theme.NodeFg,
					Background: p.theme.NodeBg,
				})
			}
		}
	}
}
