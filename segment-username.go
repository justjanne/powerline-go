package main

import (
	"fmt"
	"os"
)

func segmentUser(p *powerline) {
	var userPrompt string
	if *p.args.Shell == "bash" {
		userPrompt = " \\u "
	} else if *p.args.Shell == "zsh" {
		userPrompt = " %n "
	} else {
		user, _ := os.LookupEnv("USER")
		userPrompt = fmt.Sprintf(" %s", user)
	}

	var background uint8
	if os.Getuid() == 0 {
		background = p.theme.UsernameRootBg
	} else {
		background = p.theme.UsernameBg
	}

	p.appendSegment(segment{
		content:    userPrompt,
		foreground: p.theme.UsernameFg,
		background: background,
	})
}
