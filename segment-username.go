package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentUser(p *powerline) []pwl.Segment {
	var userPrompt string
	switch p.cfg.Shell {
	case "bash":
		userPrompt = "\\u"
	case "zsh":
		userPrompt = "%n"
	default:
		userPrompt = p.username
	}

	var background uint8
	if p.userIsAdmin {
		background = p.theme.UsernameRootBg
	} else {
		background = p.theme.UsernameBg
	}

	return []pwl.Segment{{
		Name:       "user",
		Content:    userPrompt,
		Foreground: p.theme.UsernameFg,
		Background: background,
	}}
}
