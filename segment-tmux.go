package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
)

type tmuxState struct {
	clients  int
	sessions int
	attached int
	p        *powerline
}

func (ts *tmuxState) isOn() bool {
	return ts.clients+ts.sessions+ts.attached > 0
}

func (ts *tmuxState) status() string {
	var tOn string = fmt.Sprintf("%s %s", ts.p.symbols.TmuxIndicator, ts.p.symbols.TmuxOn)
	var tOff string = fmt.Sprintf("%s %s", ts.p.symbols.TmuxIndicator, ts.p.symbols.TmuxOff)

	if !ts.p.cfg.SimpleTmux {
		tOn = fmt.Sprintf("%s c:%2d s:%2d a:%2d", ts.p.symbols.TmuxIndicator, ts.clients, ts.sessions, ts.attached)
		tOff = ""
	}

	if ts.isOn() {
		return tOn
	}

	return tOff
}

func runTmuxCmd(parameters ...string) (string, error) {
	output, err := exec.Command("tmux", parameters...).Output()

	if err != nil {
		return "", err
	}

	return string(output), nil
}

func countLines(input string, countSum bool) (int, int) {
	var cntr, sum int

	lines := strings.Split(input, "\n")

	for _, l := range lines {

		if len(l) == 0 {
			continue
		}

		cntr++

		if countSum {
			val, _ := strconv.Atoi(l)
			sum = +val
		}
	}

	return cntr, sum
}

func updateClientNum(ts *tmuxState) {
	output, err := runTmuxCmd("list-clients", "-F", "#{session_windows}")

	if err == nil {
		ts.clients, _ = countLines(output, false)
	}
}

func updateSessionNum(ts *tmuxState) {
	output, err := runTmuxCmd("list-sessions", "-F", "#{session_attached}")

	if err == nil {
		ts.sessions, ts.attached = countLines(output, true)
	}
}

func segmentTmux(p *powerline) []pwl.Segment {
	ts := &tmuxState{p: p}

	if _, exists := os.LookupEnv("TMUX"); !exists {
		updateClientNum(ts)
		updateSessionNum(ts)

		if tc := ts.status(); len(tc) > 0 {
			return []pwl.Segment{{
				Name:       "tmux",
				Content:    tc,
				Foreground: p.theme.TmuxFg,
				Background: p.theme.TmuxBg,
			}}
		}
	}

	return []pwl.Segment{}
}
