package main

func segmentRoot(p *powerline) {
	var foreground, background uint8
	if *p.args.PrevError == 0 {
		foreground = p.theme.CmdPassedFg
		background = p.theme.CmdPassedBg
	} else {
		foreground = p.theme.CmdFailedFg
		background = p.theme.CmdFailedBg
	}

	p.appendSegment(segment{
		content:    p.rootIndicator,
		foreground: foreground,
		background: background,
	})
}
