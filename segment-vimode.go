package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentViMode(p *powerline) []pwl.Segment {
	mode := p.cfg.ViMode
	if mode == "" {
		warn("'--vi-mode' is not set.")
		return []pwl.Segment{}
	}

	switch mode {
	case "vicmd":
		return []pwl.Segment{{
			Name:       "vi-mode",
			Content:    "C",
			Foreground: p.theme.ViModeCommandFg,
			Background: p.theme.ViModeCommandBg,
		}}
	default: // usually "viins" or "main"
		return []pwl.Segment{{
			Name:       "vi-mode",
			Content:    "I",
			Foreground: p.theme.ViModeInsertFg,
			Background: p.theme.ViModeInsertBg,
		}}
	}
}
