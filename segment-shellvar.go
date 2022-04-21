package main

import (
	"os"

	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentShellVar(p *powerline) []pwl.Segment {
	shellVarName := p.cfg.ShellVar
	varContent, varExists := os.LookupEnv(shellVarName)

	if !varExists {
		if shellVarName != "" {
			warn("Shell variable " + shellVarName + " does not exist.")
		}
		return []pwl.Segment{}
	}

	if varContent == "" {
		if !p.cfg.ShellVarNoWarnEmpty {
			warn("Shell variable " + shellVarName + " is empty.")
		}
		return []pwl.Segment{}
	}

	return []pwl.Segment{{
		Name:       "shell-var",
		Content:    varContent,
		Foreground: p.theme.ShellVarFg,
		Background: p.theme.ShellVarBg,
	}}
}
