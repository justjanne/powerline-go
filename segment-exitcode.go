package main

import (
	"fmt"
	"github.com/justjanne/powerline-go/exitcode"
	pwl "github.com/justjanne/powerline-go/powerline"
	"strconv"
)

var exitCodes = map[int]string{
	1:   "ERROR",
	2:   "USAGE",
	127: "NOTFOUND",
}

func getMeaningFromExitCode(exitCode int) string {
	if exitCode < 128 {
		name, ok := exitCodes[exitCode]
		if ok {
			return name
		}
	} else {
		name, ok := exitcode.Signals[exitCode-128]
		if ok {
			return name
		}
	}

	return fmt.Sprintf("%d", exitCode)
}

func segmentExitCode(p *powerline) []pwl.Segment {
	var meaning string
	if p.cfg.PrevError != 0 {
		if p.cfg.NumericExitCodes {
			meaning = strconv.Itoa(p.cfg.PrevError)
		} else {
			meaning = getMeaningFromExitCode(p.cfg.PrevError)
		}
		return []pwl.Segment{{
			Name:       "exit",
			Content:    meaning,
			Foreground: p.theme.CmdFailedFg,
			Background: p.theme.CmdFailedBg,
		}}
	}
	return []pwl.Segment{}
}
