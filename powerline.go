package main

import (
	"bytes"
	"fmt"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"
	"sync"

	pwl "github.com/justjanne/powerline-go/powerline"
	"github.com/mattn/go-runewidth"
	"github.com/shirou/gopsutil/v3/process"
	"golang.org/x/term"
	"golang.org/x/text/width"
)

// ShellInfo holds the shell information
type ShellInfo struct {
	RootIndicator         string
	ColorTemplate         string
	EscapedDollar         string
	EscapedBacktick       string
	EscapedBackslash      string
	EvalPromptPrefix      string
	EvalPromptSuffix      string
	EvalPromptRightPrefix string
	EvalPromptRightSuffix string
}

type powerline struct {
	cfg            Config
	cwd            string
	userInfo       user.User
	userIsAdmin    bool
	hostname       string
	username       string
	theme          Theme
	shell          ShellInfo
	reset          string
	symbols        SymbolTemplate
	priorities     map[string]int
	ignoreRepos    map[string]bool
	Segments       [][]pwl.Segment
	curSegment     int
	align          alignment
	rightPowerline *powerline
}

type prioritizedSegments struct {
	i    int
	segs []pwl.Segment
}

func newPowerline(cfg Config, cwd string, align alignment) *powerline {
	p := new(powerline)
	p.cfg = cfg
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
	if cfg.TrimADDomain {
		usernameWithAd := strings.SplitN(p.username, `\`, 2)
		if len(usernameWithAd) > 1 {
			// remove the Domain name from username
			p.username = usernameWithAd[1]
		}
	}
	p.userIsAdmin = userIsAdmin()

	p.theme = cfg.Themes[cfg.Theme]
	if cfg.Shell == "autodetect" {
		var shellExe string
		proc, err := process.NewProcess(int32(os.Getppid()))
		if err == nil {
			shellExe, _ = proc.Exe()
		}
		if shellExe == "" {
			shellExe = os.Getenv("SHELL")
		}
		cfg.Shell = detectShell(shellExe)
	}
	p.shell = cfg.Shells[cfg.Shell]
	p.reset = fmt.Sprintf(p.shell.ColorTemplate, "[0m")
	p.symbols = cfg.Modes[cfg.Mode]
	p.priorities = make(map[string]int)
	for idx, priority := range cfg.Priority {
		p.priorities[priority] = len(cfg.Priority) - idx
	}
	p.align = align
	p.ignoreRepos = make(map[string]bool)
	for _, r := range cfg.IgnoreRepos {
		if r == "" {
			continue
		}
		p.ignoreRepos[r] = true
	}
	p.Segments = make([][]pwl.Segment, 1)
	var mods []string
	if p.align == alignLeft {
		mods = cfg.Modules
		if len(cfg.ModulesRight) > 0 {
			if p.supportsRightModules() {
				p.rightPowerline = newPowerline(cfg, cwd, alignRight)
			} else {
				mods = append(mods, cfg.ModulesRight...)
			}
		}
	} else {
		mods = cfg.ModulesRight
	}
	initSegments(p, mods)

	return p
}

func detectShell(shellExe string) string {
	var shell string
	shellExe = path.Base(shellExe)
	if strings.Contains(shellExe, "bash") {
		shell = "bash"
	} else if strings.Contains(shellExe, "zsh") {
		shell = "zsh"
	} else {
		shell = "bare"
	}
	return shell
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
	return fmt.Sprintf(p.shell.ColorTemplate, fmt.Sprintf("[%s;5;%dm", prefix, code))
}

func (p *powerline) fgColor(code uint8) string {
	if p.theme.BoldForeground {
		return p.color("1;38", code)
	} else {
		return p.color("38", code)
	}
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
			segment.Separator = p.symbols.SeparatorReverse
		} else {
			segment.Separator = p.symbols.Separator
		}
	}
	if segment.SeparatorForeground == 0 {
		segment.SeparatorForeground = segment.Background
	}
	segment.Priority += p.priorities[origin]
	segment.Width = segment.ComputeWidth(p.cfg.Condensed)
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
	termWidth, _, err := term.GetSize(int(os.Stdin.Fd()))
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

	shellMaxLength := termWidth() * p.cfg.MaxWidthPercentage / 100
	row := p.Segments[rowNum]
	rowLength := 0

	if shellMaxLength > 0 {
		for _, segment := range row {
			rowLength += segment.Width
		}

		if rowLength > shellMaxLength && p.cfg.TruncateSegmentWidth > 0 {
			minPriorityNotTruncated := MaxInteger
			minPriorityNotTruncatedSegmentID := -1
			for idx, segment := range row {
				if segment.Width > p.cfg.TruncateSegmentWidth && segment.Priority < minPriorityNotTruncated {
					minPriorityNotTruncated = segment.Priority
					minPriorityNotTruncatedSegmentID = idx
				}
			}
			for minPriorityNotTruncatedSegmentID != -1 && rowLength > shellMaxLength {
				segment := row[minPriorityNotTruncatedSegmentID]

				rowLength -= segment.Width

				segment.Content = runewidth.Truncate(segment.Content, p.cfg.TruncateSegmentWidth-runewidth.StringWidth(segment.Separator)-3, "…")
				segment.Width = segment.ComputeWidth(p.cfg.Condensed)

				row = append(append(row[:minPriorityNotTruncatedSegmentID], segment), row[minPriorityNotTruncatedSegmentID+1:]...)
				rowLength += segment.Width

				minPriorityNotTruncated = MaxInteger
				minPriorityNotTruncatedSegmentID = -1
				for idx, segment := range row {
					if segment.Width > p.cfg.TruncateSegmentWidth && segment.Priority < minPriorityNotTruncated {
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
	if !p.cfg.EastAsianWidth {
		return 0
	}
	numEastAsianRunes := 0
	for _, r := range *segmentContent {
		switch width.LookupRune(r).Kind() {
		case width.EastAsianAmbiguous:
			numEastAsianRunes++
		case width.Neutral:
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
		if !p.cfg.Condensed {
			buffer.WriteRune(' ')
		}
		buffer.WriteString(segment.Content)
		numEastAsianRunes += p.numEastAsianRunes(&segment.Content)
		if !p.cfg.Condensed {
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

	if p.cfg.Eval {
		if p.align == alignLeft {
			buffer.WriteString(p.shell.EvalPromptPrefix)
		} else if p.supportsRightModules() {
			buffer.WriteString(p.shell.EvalPromptRightPrefix)
		}
	}

	for rowNum := range p.Segments {
		p.truncateRow(rowNum)
		p.drawRow(rowNum, &buffer)
		if rowNum < len(p.Segments)-1 {
			buffer.WriteRune('\n')
		}
	}

	if p.cfg.PromptOnNewLine {
		buffer.WriteRune('\n')

		var foreground, background uint8
		if p.cfg.PrevError == 0 || p.cfg.StaticPromptIndicator {
			foreground = p.theme.CmdPassedFg
			background = p.theme.CmdPassedBg
		} else {
			foreground = p.theme.CmdFailedFg
			background = p.theme.CmdFailedBg
		}

		buffer.WriteString(p.fgColor(foreground))
		buffer.WriteString(p.bgColor(background))
		buffer.WriteString(p.shell.RootIndicator)
		buffer.WriteString(p.reset)
		buffer.WriteString(p.fgColor(background))
		buffer.WriteString(p.symbols.Separator)
		buffer.WriteString(p.reset)
		buffer.WriteRune(' ')
	}

	if p.cfg.Eval {
		switch p.align {
		case alignLeft:
			buffer.WriteString(p.shell.EvalPromptSuffix)
			if p.supportsRightModules() {
				buffer.WriteRune('\n')
				if !p.hasRightModules() {
					buffer.WriteString(p.shell.EvalPromptRightPrefix + p.shell.EvalPromptRightSuffix)
				}
			}
		case alignRight:
			if p.supportsRightModules() {
				buffer.Truncate(buffer.Len() - 1)
				buffer.WriteString(p.shell.EvalPromptRightSuffix)
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
	return p.shell.EvalPromptRightPrefix != "" || p.shell.EvalPromptRightSuffix != ""
}

func (p *powerline) isRightPrompt() bool {
	return p.align == alignRight && p.supportsRightModules()
}
