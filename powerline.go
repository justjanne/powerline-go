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
	pathAliases     map[string]string
	theme           Theme
	shellInfo       ShellInfo
	reset           string
	symbolTemplates Symbols
	priorities      map[string]int
	ignoreRepos     map[string]bool
	Segments        [][]segment
	curSegment      int
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
		if r == "" {
			continue
		}
		p.ignoreRepos[r] = true
	}
	p.pathAliases = make(map[string]string)
	for _, pa := range strings.Split(*args.PathAliases, ",") {
		if pa == "" {
			continue
		}
		kv := strings.SplitN(pa, "=", 2)
		p.pathAliases[kv[0]] = kv[1]
	}
	p.Segments = make([][]segment, 1)
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
	p.Segments[p.curSegment] = append(p.Segments[p.curSegment], segment)
}

func (p *powerline) newRow() {
	p.Segments = append(p.Segments, make([]segment, 0))
	p.curSegment = p.curSegment + 1
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

func (p *powerline) truncateRow(rowNum int) {

	shellMaxLength := termWidth() * *p.args.MaxWidthPercentage / 100
	row := p.Segments[rowNum]
	rowLength := 0

	if shellMaxLength > 0 {
		for _, segment := range row {
			rowLength += segment.width
		}

		if rowLength > shellMaxLength && *p.args.TruncateSegmentWidth > 0 {
			minPriorityNotTruncated := MaxInteger
			minPriorityNotTruncatedSegmentId := -1
			for idx, segment := range row {
				if segment.width > *p.args.TruncateSegmentWidth && segment.priority < minPriorityNotTruncated {
					minPriorityNotTruncated = segment.priority
					minPriorityNotTruncatedSegmentId = idx
				}
			}
			for minPriorityNotTruncatedSegmentId != -1 && rowLength > shellMaxLength {
				segment := row[minPriorityNotTruncatedSegmentId]

				rowLength -= segment.width

				segment.content = runewidth.Truncate(segment.content, *p.args.TruncateSegmentWidth-runewidth.StringWidth(segment.separator)-3, "…")
				segment.width = segment.computeWidth()

				row = append(append(row[:minPriorityNotTruncatedSegmentId], segment), row[minPriorityNotTruncatedSegmentId+1:]...)
				rowLength += segment.width

				minPriorityNotTruncated = MaxInteger
				minPriorityNotTruncatedSegmentId = -1
				for idx, segment := range row {
					if segment.width > *p.args.TruncateSegmentWidth && segment.priority < minPriorityNotTruncated {
						minPriorityNotTruncated = segment.priority
						minPriorityNotTruncatedSegmentId = idx
					}
				}
			}
		}

		for rowLength > shellMaxLength {
			minPriority := MaxInteger
			minPrioritySegmentId := -1
			for idx, segment := range row {
				if segment.priority < minPriority {
					minPriority = segment.priority
					minPrioritySegmentId = idx
				}
			}
			if minPrioritySegmentId != -1 {
				segment := row[minPrioritySegmentId]
				row = append(row[:minPrioritySegmentId], row[minPrioritySegmentId+1:]...)
				rowLength -= segment.width
			}
		}
	}
	p.Segments[rowNum] = row
}

func (p *powerline) numEastAsianRunes(segmentContent *string) int {
	if *p.args.EastAsianWidth {
		return 0
	}
	numEastAsianRunes := 0
	for _, r := range *segmentContent {
		switch width.LookupRune(r).Kind() {
		case width.Neutral:
		case width.EastAsianAmbiguous:
			numEastAsianRunes += 1
		case width.EastAsianWide:
		case width.EastAsianNarrow:
		case width.EastAsianFullwidth:
		case width.EastAsianHalfwidth:
		}
	}
	return numEastAsianRunes
}

func (p *powerline) drawRow(rowNum int, buffer *bytes.Buffer) {
	row := p.Segments[rowNum]
	numEastAsianRunes := 0
	for idx, segment := range row {
        if (segment.hideSeparators) {
            buffer.WriteString(segment.content);
            continue;
        }
		var separatorBackground string
		if idx >= len(row)-1 {
			separatorBackground = p.reset
		} else {
			nextSegment := row[idx+1]
			separatorBackground = p.bgColor(nextSegment.background)
		}
		buffer.WriteString(p.fgColor(segment.foreground))
		buffer.WriteString(p.bgColor(segment.background))
		buffer.WriteRune(' ')
		buffer.WriteString(segment.content)
		numEastAsianRunes += p.numEastAsianRunes(&segment.content)
		buffer.WriteRune(' ')
		buffer.WriteString(separatorBackground)
		buffer.WriteString(p.fgColor(segment.separatorForeground))
		buffer.WriteString(segment.separator)
		buffer.WriteString(p.reset)
	}
	buffer.WriteRune(' ')

	for i := 0; i < numEastAsianRunes; i++ {
		buffer.WriteRune(' ')
	}
}

func (p *powerline) draw() string {

	var buffer bytes.Buffer

	for rowNum := range p.Segments {
		p.truncateRow(rowNum)
		p.drawRow(rowNum, &buffer)
		if rowNum < len(p.Segments)-1 {
			buffer.WriteRune('\n')
		}
	}

	if *p.args.PromptOnNewLine {
		buffer.WriteRune('\n')

		var foreground, background uint8
		if *p.args.PrevError == 0 {
			foreground = p.theme.CmdPassedFg
			background = p.theme.CmdPassedBg
		} else {
			foreground = p.theme.CmdFailedFg
			background = p.theme.CmdFailedBg
		}

		buffer.WriteString(p.fgColor(foreground))
		buffer.WriteString(p.bgColor(background))
		buffer.WriteString(p.shellInfo.rootIndicator)
		buffer.WriteString(p.reset)
		buffer.WriteString(p.fgColor(background))
		buffer.WriteString(p.symbolTemplates.Separator)
		buffer.WriteString(p.reset)
		buffer.WriteRune(' ')
	}

	return buffer.String()
}
