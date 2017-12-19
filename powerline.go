package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/text/width"
	"os"
	"strconv"
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
	priorities      map[string]int
	ignoreRepos     map[string]bool
	Segments        []segment
}

func NewPowerline(args args, cwd string, priorities map[string]int) *powerline {
	p := new(powerline)
	p.args = args
	p.cwd = cwd
	p.theme = themes[*args.Theme]
	p.shellInfo = shellInfos[*args.Shell]
	p.reset = fmt.Sprintf(p.shellInfo.colorTemplate, "[0m")
	p.symbolTemplates = symbolTemplates[*args.Mode]
	p.priorities = priorities
	p.ignoreRepos = make(map[string]bool)
	for _, r := range strings.Split(*args.IgnoreRepos, ",") {
		p.ignoreRepos[r] = true
	}
	p.Segments = make([]segment, 0)
	return p
}

func (p *powerline) color(prefix string, code uint8) string {
	if code == p.theme.Reset {
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

func (p *powerline) appendSegment(origin string, segment segment) {
	if segment.separator == "" {
		segment.separator = p.symbolTemplates.Separator
	}
	if segment.separatorForeground == 0 {
		segment.separatorForeground = segment.background
	}
	priority, _ := p.priorities[origin]
	segment.priority += priority
	segment.width = segment.computeWidth()
	p.Segments = append(p.Segments, segment)
}

func termWidth() int {
	termWidth, _, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		shellMaxLengthStr, found := os.LookupEnv("COLUMNS")
		if !found {
			return 0
		}

		shellMaxLength64, err := strconv.ParseInt(shellMaxLengthStr, 0, 64)
		if err != nil {
			return 0
		}

		termWidth = int(shellMaxLength64)
	}

	return termWidth
}

func (p *powerline) draw() string {
	shellMaxLength := termWidth()

	shellMaxLength = shellMaxLength * *p.args.MaxWidthPercentage / 100

	shellActualLength := 0
	if shellMaxLength > 0 {
		for _, segment := range p.Segments {
			shellActualLength += segment.width
		}
		if shellActualLength > shellMaxLength && *p.args.TruncateSegmentWidth > 0 {
			minPriorityNotTruncated := MaxInteger
			minPriorityNotTruncatedSegmentId := -1
			for idx, segment := range p.Segments {
				if segment.width > *p.args.TruncateSegmentWidth && segment.priority < minPriorityNotTruncated {
					minPriorityNotTruncated = segment.priority
					minPriorityNotTruncatedSegmentId = idx
				}
			}
			for minPriorityNotTruncatedSegmentId != -1 && shellActualLength > shellMaxLength {
				segment := p.Segments[minPriorityNotTruncatedSegmentId]

				shellActualLength -= segment.width

				segment.content = runewidth.Truncate(segment.content, *p.args.TruncateSegmentWidth-runewidth.StringWidth(segment.separator)-3, "…")
				segment.width = segment.computeWidth()

				p.Segments = append(append(p.Segments[:minPriorityNotTruncatedSegmentId], segment), p.Segments[minPriorityNotTruncatedSegmentId+1:]...)
				shellActualLength += segment.width

				minPriorityNotTruncated = MaxInteger
				minPriorityNotTruncatedSegmentId = -1
				for idx, segment := range p.Segments {
					if segment.width > *p.args.TruncateSegmentWidth && segment.priority < minPriorityNotTruncated {
						minPriorityNotTruncated = segment.priority
						minPriorityNotTruncatedSegmentId = idx
					}
				}
			}
		}

		for shellActualLength > shellMaxLength {
			minPriority := MaxInteger
			minPrioritySegmentId := -1
			for idx, segment := range p.Segments {
				if segment.priority < minPriority {
					minPriority = segment.priority
					minPrioritySegmentId = idx
				}
			}
			if minPrioritySegmentId != -1 {
				segment := p.Segments[minPrioritySegmentId]
				p.Segments = append(p.Segments[:minPrioritySegmentId], p.Segments[minPrioritySegmentId+1:]...)
				shellActualLength -= segment.width
			}
		}
	}

	var buffer bytes.Buffer
	for idx, segment := range p.Segments {
		if (segment.special) {
			buffer.WriteString(segment.content);
			continue;
		}

		var separatorBackground string
		if idx >= len(p.Segments)-1 {
			separatorBackground = p.reset
		} else {
			nextSegment := p.Segments[idx+1]
			separatorBackground = p.bgColor(nextSegment.background)
		}

		buffer.WriteString(p.fgColor(segment.foreground))
		buffer.WriteString(p.bgColor(segment.background))
		buffer.WriteRune(' ')
		buffer.WriteString(segment.content)
		buffer.WriteRune(' ')
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
