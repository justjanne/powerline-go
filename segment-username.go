package main

import (
	"os"
	"os/user"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentUser(p *powerline) {
	var userPrompt string
	if *p.args.Shell == "bash" {
		userPrompt = "\\u"
	} else if *p.args.Shell == "zsh" {
		userPrompt = "%n"
	} else {
		if userName, found := os.LookupEnv("USER"); found {
			userPrompt = userName
		} else {
			userInfo, err := user.Current()
			if err == nil {
				userPrompt = userInfo.Username
			}
		}
	}

	var background uint8
	if os.Getuid() == 0 {
		background = p.theme.UsernameRootBg
	} else {
		background = p.theme.UsernameBg
	}

	p.appendSegment("user", pwl.Segment{
		Content:    userPrompt,
		Foreground: p.theme.UsernameFg,
		Background: background,
	})
}
