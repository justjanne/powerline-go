package main

import (
	"bytes"
	"fmt"
	"golang.org/x/text/width"
)

type ShellInfo struct {
	rootIndicator    string
	colorTemplate    string
	escapedDollar    string
	escapedBacktick  string
	escapedBackslash string
}

type powerline struct {
	args            args
	cwd             string
	theme           Theme
	shellInfo       ShellInfo
	reset           string
	symbolTemplates Symbols
	Segments        []segment
}

func NewPowerline(args args, cwd string, theme Theme) *powerline {
	p := new(powerline)
	p.args = args
	p.cwd = cwd
	p.theme = theme
	p.shellInfo = shellInfos[*args.Shell]
	p.reset = fmt.Sprintf(p.shellInfo.colorTemplate, "[0m")
	p.symbolTemplates = symbolTemplates[*args.Mode]
	p.Segments = make([]segment, 0)
	return p
}

func (p *powerline) color(prefix string, code uint8) string {
	if code == defaultTheme.Reset {
		return p.reset
	} else {
		return fmt.Sprintf(p.shellInfo.colorTemplate, fmt.Sprintf("[%s;5;%dm", prefix, code))
	}
}

func (p *powerline) fgColor(code uint8) string {
	return p.color("38", code)
}

func (p *powerline) bgColor(code uint8) string {
	return p.color("48", code)
}

func (p *powerline) appendSegment(segment segment) {
	if segment.separator == "" {
		segment.separator = p.symbolTemplates.Separator
	}
	if segment.separatorForeground == 0 {
		segment.separatorForeground = segment.background
	}
	p.Segments = append(p.Segments, segment)
}

func (p *powerline) draw() string {
	var buffer bytes.Buffer
	for idx := range p.Segments {
		segment := p.Segments[idx]

		var separatorBackground string
		if idx >= len(p.Segments)-1 {
			separatorBackground = p.reset
		} else {
			nextSegment := p.Segments[idx+1]
			separatorBackground = p.bgColor(nextSegment.background)
		}

		buffer.WriteString(p.fgColor(segment.foreground))
		buffer.WriteString(p.bgColor(segment.background))
		buffer.WriteString(segment.content)
		buffer.WriteString(separatorBackground)
		buffer.WriteString(p.fgColor(segment.separatorForeground))
		buffer.WriteString(segment.separator)
		buffer.WriteString(p.reset)
	}
	buffer.WriteRune(' ')

	drawnResult := buffer.String()
	if *p.args.EastAsianWidth {
		var spaceBuffer bytes.Buffer
		for _, r := range drawnResult {
			switch width.LookupRune(r).Kind() {
			case width.Neutral:
			case width.EastAsianAmbiguous:
				spaceBuffer.WriteRune(' ')
			case width.EastAsianWide:
			case width.EastAsianNarrow:
			case width.EastAsianFullwidth:
			case width.EastAsianHalfwidth:
			}
		}
		drawnResult += spaceBuffer.String()
	}

	if *p.args.PromptOnNewLine {
		var nextLineBuffer bytes.Buffer
		nextLineBuffer.WriteRune('\n')

		var foreground, background uint8
		if *p.args.PrevError == 0 {
			foreground = p.theme.CmdPassedFg
			background = p.theme.CmdPassedBg
		} else {
			foreground = p.theme.CmdFailedFg
			background = p.theme.CmdFailedBg
		}

		nextLineBuffer.WriteString(p.fgColor(foreground))
		nextLineBuffer.WriteString(p.bgColor(background))
		nextLineBuffer.WriteString(p.shellInfo.rootIndicator)
		nextLineBuffer.WriteString(p.reset)
		nextLineBuffer.WriteString(p.fgColor(background))
		nextLineBuffer.WriteString(p.symbolTemplates.Separator)
		nextLineBuffer.WriteString(p.reset)
		nextLineBuffer.WriteRune(' ')

		drawnResult += nextLineBuffer.String()
	}

	return drawnResult
}
