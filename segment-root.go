package main

import pwl "github.com/justjanne/powerline-go/powerline"

func segmentRoot(p *powerline) {
	var foreground, background uint8
	if *p.args.PrevError == 0 {
		foreground = p.theme.CmdPassedFg
		background = p.theme.CmdPassedBg
	} else {
		foreground = p.theme.CmdFailedFg
		background = p.theme.CmdFailedBg
	}

	p.appendSegment("root", pwl.Segment{
		Content:    p.shellInfo.rootIndicator,
		Foreground: foreground,
		Background: background,
	})
}
