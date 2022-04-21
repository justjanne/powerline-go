package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
)

const pkgfile = "./package.json"

type packageJSON struct {
	Version string `json:"version"`
}

func getNodeVersion() string {
	out, err := exec.Command("node", "--version").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSuffix(string(out), "\n")
}

func getPackageVersion() string {
	stat, err := os.Stat(pkgfile)
	if err != nil {
		return ""
	}
	if stat.IsDir() {
		return ""
	}
	pkg := packageJSON{""}
	raw, err := ioutil.ReadFile(pkgfile)
	if err != nil {
		return ""
	}
	err = json.Unmarshal(raw, &pkg)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(pkg.Version)
}

func segmentNode(p *powerline) []pwl.Segment {
	nodeVersion := getNodeVersion()
	packageVersion := getPackageVersion()

	segments := []pwl.Segment{}

	if nodeVersion != "" {
		segments = append(segments, pwl.Segment{
			Name:       "node",
			Content:    p.symbols.NodeIndicator + " " + nodeVersion,
			Foreground: p.theme.NodeVersionFg,
			Background: p.theme.NodeVersionBg,
		})
	}

	if packageVersion != "" {
		segments = append(segments, pwl.Segment{
			Name:       "node-segment",
			Content:    packageVersion + " " + p.symbols.NodeIndicator,
			Foreground: p.theme.NodeFg,
			Background: p.theme.NodeBg,
		})
	}

	return segments
}
