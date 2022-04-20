package main

import (
	"os"
	"path/filepath"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentDirenv(p *powerline) []pwl.Segment {
	content := os.Getenv("DIRENV_DIR")
	if content == "" {
		return []pwl.Segment{}
	}
	if strings.TrimPrefix(content, "-") == p.userInfo.HomeDir {
		content = "~"
	} else {
		content = filepath.Base(content)
	}

	return []pwl.Segment{{
		Name:       "direnv",
		Content:    content,
		Foreground: p.theme.DotEnvFg,
		Background: p.theme.DotEnvBg,
	}}
}
