package main

import (
	"os"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentScreenName(p *powerline) []pwl.Segment {
	var env string
	if env == "" {
		env, _ = os.LookupEnv("STY")
	}
	if env == "" {
		return []pwl.Segment{}
	}
	// envName := env
	envName := strings.Split(env, ".")[1]
	if p.cfg.VenvNameSizeLimit > 0 && len(envName) > p.cfg.VenvNameSizeLimit {
		envName = p.symbols.VenvIndicator
	}

	return []pwl.Segment{{
		Name:       "screen",
		Content:    envName,
		Foreground: p.theme.ScreenFg,
		Background: p.theme.ScreenBg,
	}}
}
