package main

import pwl "github.com/justjanne/powerline-go/powerline"

func segmentRoot(p *powerline) []pwl.Segment {
	var foreground, background uint8
	if *p.args.PrevError == 0 || *p.args.StaticPromptIndicator {
		foreground = p.theme.CmdPassedFg
		background = p.theme.CmdPassedBg
	} else {
		foreground = p.theme.CmdFailedFg
		background = p.theme.CmdFailedBg
	}

	return []pwl.Segment{{
		Name:       "root",
		Content:    p.shellInfo.rootIndicator,
		Foreground: foreground,
		Background: background,
	}}
}
