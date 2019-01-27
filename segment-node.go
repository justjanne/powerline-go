package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
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
				p.appendSegment("node-version", segment{
					content:    pkg.Version + " \u2B22",
					foreground: p.theme.NodeFg,
					background: p.theme.NodeBg,
				})
			}
		}
	}
}
