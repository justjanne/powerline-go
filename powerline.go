package main

import (
	"bytes"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"sync"

	pwl "github.com/justjanne/powerline-go/powerline"
	"github.com/mattn/go-runewidth"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/text/width"
)

// ShellInfo holds the shell information
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
	userInfo               user.User
	hostname               string
	username               string
	pathAliases            map[string]string
	theme                  Theme
	shellInfo              ShellInfo
	reset                  string
	symbolTemplates        Symbols
	priorities             map[string]int
	ignoreRepos            map[string]bool
	Segments               [][]pwl.Segment
	curSegment             int
	align                  alignment
	rightPowerline         *powerline
	appendEastAsianPadding int
}

type prioritizedSegments struct {
	i    int
	segs []pwl.Segment
}

func newPowerline(args args, cwd string, priorities map[string]int, align alignment) *powerline {
	p := new(powerline)
	p.args = args
	p.cwd = cwd
	userInfo, err := user.Current()
	if userInfo != nil && err == nil {
		p.userInfo = *userInfo
	}
	p.hostname, _ = os.Hostname()

	hostnamePrefix := fmt.Sprintf("%s%c", p.hostname, os.PathSeparator)
	if strings.HasPrefix(p.userInfo.Username, hostnamePrefix) {
		p.username = p.userInfo.Username[len(hostnamePrefix):]
	} else {
		p.username = p.userInfo.Username
	}
	if args.TrimADDomain != nil && *args.TrimADDomain {
		usernameWithAd := strings.SplitN(p.username, `\`, 2)
		if len(usernameWithAd) > 1 {
			// remove the Domain name from username
			p.username = usernameWithAd[1]
		}
	}

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
	p.Segments = make([][]pwl.Segment, 1)
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
	initSegments(p, strings.Split(mods, ","))

	return p
}

func initSegments(p *powerline, mods []string) {
	orderedSegments := map[int][]pwl.Segment{}
	c := make(chan prioritizedSegments, len(mods))
	wg := sync.WaitGroup{}
	for i, module := range mods {
		wg.Add(1)
		go func(w *sync.WaitGroup, i int, module string, c chan prioritizedSegments) {
			elem, ok := modules[module]
			if ok {
				c <- prioritizedSegments{
					i:    i,
					segs: elem(p),
				}
			} else {
				s, ok := segmentPlugin(p, module)
				if ok {
					c <- prioritizedSegments{
						i:    i,
						segs: s,
					}
				} else {
					println("Module not found: " + module)
				}
			}
			wg.Done()
		}(&wg, i, module, c)
	}
	wg.Wait()
	close(c)
	for s := range c {
		orderedSegments[s.i] = s.segs
	}
	for i := 0; i < len(mods); i++ {
		for _, seg := range orderedSegments[i] {
			p.appendSegment(seg.Name, seg)
		}
	}
}

func (p *powerline) color(prefix string, code uint8) string {
	if code == p.theme.Reset {
		return p.reset
	}
	return fmt.Sprintf(p.shellInfo.colorTemplate, fmt.Sprintf("[%s;5;%dm", prefix, code))
}

func (p *powerline) fgColor(code uint8) string {
	return p.color("38", code)
}

func (p *powerline) bgColor(code uint8) string {
	return p.color("48", code)
}

func (p *powerline) appendSegment(origin string, segment pwl.Segment) {
	if segment.Foreground == segment.Background && segment.Background == 0 {
		segment.Background = p.theme.DefaultBg
		segment.Foreground = p.theme.DefaultFg
	}
	if segment.Separator == "" {
		if p.isRightPrompt() {
			segment.Separator = p.symbolTemplates.SeparatorReverse
		} else {
			segment.Separator = p.symbolTemplates.Separator
		}
	}
	if segment.SeparatorForeground == 0 {
		segment.SeparatorForeground = segment.Background
	}
	priority, _ := p.priorities[origin]
	segment.Priority += priority
	segment.Width = segment.ComputeWidth(*p.args.Condensed)
	if segment.NewLine {
		p.newRow()
	} else {
		p.Segments[p.curSegment] = append(p.Segments[p.curSegment], segment)
	}
}

func (p *powerline) newRow() {
	if len(p.Segments[p.curSegment]) > 0 {
		p.Segments = append(p.Segments, make([]pwl.Segment, 0))
		p.curSegment = p.curSegment + 1
	}
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
			rowLength += segment.Width
		}

		if rowLength > shellMaxLength && *p.args.TruncateSegmentWidth > 0 {
			minPriorityNotTruncated := MaxInteger
			minPriorityNotTruncatedSegmentID := -1
			for idx, segment := range row {
				if segment.Width > *p.args.TruncateSegmentWidth && segment.Priority < minPriorityNotTruncated {
					minPriorityNotTruncated = segment.Priority
					minPriorityNotTruncatedSegmentID = idx
				}
			}
			for minPriorityNotTruncatedSegmentID != -1 && rowLength > shellMaxLength {
				segment := row[minPriorityNotTruncatedSegmentID]

				rowLength -= segment.Width

				segment.Content = runewidth.Truncate(segment.Content, *p.args.TruncateSegmentWidth-runewidth.StringWidth(segment.Separator)-3, "…")
				segment.Width = segment.ComputeWidth(*p.args.Condensed)

				row = append(append(row[:minPriorityNotTruncatedSegmentID], segment), row[minPriorityNotTruncatedSegmentID+1:]...)
				rowLength += segment.Width

				minPriorityNotTruncated = MaxInteger
				minPriorityNotTruncatedSegmentID = -1
				for idx, segment := range row {
					if segment.Width > *p.args.TruncateSegmentWidth && segment.Priority < minPriorityNotTruncated {
						minPriorityNotTruncated = segment.Priority
						minPriorityNotTruncatedSegmentID = idx
					}
				}
			}
		}

		for rowLength > shellMaxLength {
			minPriority := MaxInteger
			minPrioritySegmentID := -1
			for idx, segment := range row {
				if segment.Priority < minPriority {
					minPriority = segment.Priority
					minPrioritySegmentID = idx
				}
			}
			if minPrioritySegmentID != -1 {
				segment := row[minPrioritySegmentID]
				row = append(row[:minPrioritySegmentID], row[minPrioritySegmentID+1:]...)
				rowLength -= segment.Width
			}
		}
	}
	p.Segments[rowNum] = row
}

func (p *powerline) numEastAsianRunes(segmentContent *string) int {
	if !*p.args.EastAsianWidth {
		return 0
	}
	numEastAsianRunes := 0
	for _, r := range *segmentContent {
		switch width.LookupRune(r).Kind() {
		case width.Neutral:
		case width.EastAsianAmbiguous:
			numEastAsianRunes++
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
		if segment.HideSeparators {
			buffer.WriteString(segment.Content)
			continue
		}
		var separatorBackground string
		if p.isRightPrompt() {
			if idx == 0 {
				separatorBackground = p.reset
			} else {
				prevSegment := row[idx-1]
				separatorBackground = p.bgColor(prevSegment.Background)
			}
			buffer.WriteString(separatorBackground)
			buffer.WriteString(p.fgColor(segment.SeparatorForeground))
			buffer.WriteString(segment.Separator)
		} else {
			if idx >= len(row)-1 {
				if !p.hasRightModules() || p.supportsRightModules() {
					separatorBackground = p.reset
				} else if p.hasRightModules() && rowNum >= len(p.Segments)-1 {
					nextSegment := p.rightPowerline.Segments[0][0]
					separatorBackground = p.bgColor(nextSegment.Background)
				}
			} else {
				nextSegment := row[idx+1]
				separatorBackground = p.bgColor(nextSegment.Background)
			}
		}
		buffer.WriteString(p.fgColor(segment.Foreground))
		buffer.WriteString(p.bgColor(segment.Background))
		if !*p.args.Condensed {
			buffer.WriteRune(' ')
		}
		buffer.WriteString(segment.Content)
		numEastAsianRunes += p.numEastAsianRunes(&segment.Content)
		if !*p.args.Condensed {
			buffer.WriteRune(' ')
		}
		if !p.isRightPrompt() {
			buffer.WriteString(separatorBackground)
			buffer.WriteString(p.fgColor(segment.SeparatorForeground))
			buffer.WriteString(segment.Separator)
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
