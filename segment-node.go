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

func segmentNode(p *powerline) []pwl.Segment {
	stat, err := os.Stat(pkgfile)
	if err == nil && !stat.IsDir() {
		pkg := packageJSON{"!"}
		raw, err := ioutil.ReadFile(pkgfile)
		if err == nil {
			err = json.Unmarshal(raw, &pkg)
			if err == nil {
				return []pwl.Segment{{
					Name:       "node-segment",
					Content:    pkg.Version + " \u2B22",
					Foreground: p.theme.NodeFg,
					Background: p.theme.NodeBg,
				}}
			}
		}
	}
	return []pwl.Segment{}
}
