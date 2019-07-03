package main

import (
	"bytes"
	"fmt"
	"strings"

	"os"
	"strconv"

	"github.com/mattn/go-runewidth"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/text/width"
)

type ShellInfo struct {
	rootIndicator         string
	colorTemplate         string
	escapedDollar         string
	escapedBacktick       string
	escapedBackslash      string
	evalPromptPrefix      string
	evalPromptSuffix      string
	evalPromptRightPrefix string
	evalPromptRightSuffix string
}

type powerline struct {
	args                   args
	cwd                    string
	pathAliases            map[string]string
	theme                  Theme
	shellInfo              ShellInfo
	reset                  string
	symbolTemplates        Symbols
	priorities             map[string]int
	ignoreRepos            map[string]bool
	Segments               [][]segment
	curSegment             int
	align                  alignment
	rightPowerline         *powerline
	appendEastAsianPadding int
}

func newPowerline(args args, cwd string, priorities map[string]int, align alignment) *powerline {
	p := new(powerline)
	p.args = args
	p.cwd = cwd
	p.theme = themes[*args.Theme]
	p.shellInfo = shellInfos[*args.Shell]
	p.reset = fmt.Sprintf(p.shellInfo.colorTemplate, "[0m")
	p.symbolTemplates = symbolTemplates[*args.Mode]
	p.priorities = priorities
	p.align = align
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
	var mods string
	if p.align == alignLeft {
		mods = *args.Modules
		if len(*args.ModulesRight) > 0 {
			if p.supportsRightModules() {
				p.rightPowerline = newPowerline(args, cwd, priorities, alignRight)
			} else {
				mods += `,` + *args.ModulesRight
			}
		}
	} else {
		mods = *args.ModulesRight
	}
	for _, module := range strings.Split(mods, ",") {
		elem, ok := modules[module]
		if !ok {
			println("Module not found: " + module)
			continue
		}
		elem(p)
	}
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
		if p.isRightPrompt() {
			segment.separator = p.symbolTemplates.SeparatorReverse
		} else {
			segment.separator = p.symbolTemplates.Separator
		}
	}
	if segment.separatorForeground == 0 {
		segment.separatorForeground = segment.background
	}
	priority, _ := p.priorities[origin]
	segment.priority += priority
	segment.width = segment.computeWidth(*p.args.Condensed)
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
				segment.width = segment.computeWidth(*p.args.Condensed)

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

	// Prepend padding
	if p.isRightPrompt() {
		buffer.WriteRune(' ')
	}
	for idx, segment := range row {
		if segment.hideSeparators {
			buffer.WriteString(segment.content)
			continue
		}
		var separatorBackground string
		if p.isRightPrompt() {
			if idx == 0 {
				separatorBackground = p.reset
			} else {
				prevSegment := row[idx-1]
				separatorBackground = p.bgColor(prevSegment.background)
			}
			buffer.WriteString(separatorBackground)
			buffer.WriteString(p.fgColor(segment.separatorForeground))
			buffer.WriteString(segment.separator)
		} else {
			if idx >= len(row)-1 {
				if !p.hasRightModules() || p.supportsRightModules() {
					separatorBackground = p.reset
				} else if p.hasRightModules() && rowNum >= len(p.Segments)-1 {
					nextSegment := p.rightPowerline.Segments[0][0]
					separatorBackground = p.bgColor(nextSegment.background)
				}
			} else {
				nextSegment := row[idx+1]
				separatorBackground = p.bgColor(nextSegment.background)
			}
		}
		buffer.WriteString(p.fgColor(segment.foreground))
		buffer.WriteString(p.bgColor(segment.background))
		if !*p.args.Condensed {
			buffer.WriteRune(' ')
		}
		buffer.WriteString(segment.content)
		numEastAsianRunes += p.numEastAsianRunes(&segment.content)
		if !*p.args.Condensed {
			buffer.WriteRune(' ')
		}
		if !p.isRightPrompt() {
			buffer.WriteString(separatorBackground)
			buffer.WriteString(p.fgColor(segment.separatorForeground))
			buffer.WriteString(segment.separator)
		}
		buffer.WriteString(p.reset)
	}

	// Append padding before cursor for left-aligned prompts
	if !p.isRightPrompt() || !p.hasRightModules() {
		buffer.WriteRune(' ')
	}

	// Don't append padding for right-aligned modules
	if !p.isRightPrompt() {
		for i := 0; i < numEastAsianRunes; i++ {
			buffer.WriteRune(' ')
		}
	}
}

func (p *powerline) draw() string {

	var buffer bytes.Buffer

	if *p.args.Eval {
		if p.align == alignLeft {
			buffer.WriteString(p.shellInfo.evalPromptPrefix)
		} else if p.supportsRightModules() {
			buffer.WriteString(p.shellInfo.evalPromptRightPrefix)
		}
	}

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
		if *p.args.PrevError == 0 || *p.args.StaticPromptIndicator {
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

	if *p.args.Eval {
		switch p.align {
		case alignLeft:
			buffer.WriteString(p.shellInfo.evalPromptSuffix)
			if p.supportsRightModules() {
				buffer.WriteRune('\n')
				if !p.hasRightModules() {
					buffer.WriteString(p.shellInfo.evalPromptRightPrefix + p.shellInfo.evalPromptRightSuffix)
				}
			}
		case alignRight:
			if p.supportsRightModules() {
				buffer.WriteString(p.shellInfo.evalPromptSuffix)
			}
		}
		if p.hasRightModules() {
			buffer.WriteString(p.rightPowerline.draw())
		}
	}

	return buffer.String()
}

func (p *powerline) hasRightModules() bool {
	return p.rightPowerline != nil && len(p.rightPowerline.Segments[0]) > 0
}

func (p *powerline) supportsRightModules() bool {
	return p.shellInfo.evalPromptRightPrefix != "" || p.shellInfo.evalPromptRightSuffix != ""
}

func (p *powerline) isRightPrompt() bool {
	return p.align == alignRight && p.supportsRightModules()
}
