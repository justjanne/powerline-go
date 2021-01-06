package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
)

func segmentShellVar(p *powerline) []pwl.Segment {
	shellVarName := *p.args.ShellVar
	varContent, varExists := os.LookupEnv(shellVarName)

	if varExists {
		if varContent != "" {
			return []pwl.Segment{{
				Name:       "shell-var",
				Content:    varContent,
				Foreground: p.theme.ShellVarFg,
				Background: p.theme.ShellVarBg,
			}}
		}
		if !*p.args.ShellVarNoWarnEmpty {
			warn("Shell variable " + shellVarName + " is empty.")
		}
	} else {
		warn("Shell variable " + shellVarName + " does not exist.")
	}
	return []pwl.Segment{}
}
