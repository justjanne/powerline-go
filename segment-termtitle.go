package main

// Port of set_term_title segment from powerine-shell:
// https://github.com/b-ryan/powerline-shell/blob/master/powerline_shell/segments/set_term_title.py

import (
	"fmt"
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
	"strings"
)

func segmentTermTitle(p *powerline) []pwl.Segment {
	var title string

	term := os.Getenv("TERM")
	if !(strings.Contains(term, "xterm") || strings.Contains(term, "rxvt")) {
		return []pwl.Segment{}
	}

	if *p.args.Shell == "bash" {
		title = "\\[\\e]0;\\u@\\h: \\w\\a\\]"
	} else if *p.args.Shell == "zsh" {
		title = "%{\033]0;%n@%m: %~\007%}"
	} else {
		user := os.Getenv("USER")
		host, _ := os.Hostname()
		cwd := p.cwd
		title = fmt.Sprintf("\033]0;%s@%s: %s\007", user, host, cwd)
	}

	return []pwl.Segment{{
		Name:           "termtitle",
		Content:        title,
		Priority:       MaxInteger, // do not truncate
		HideSeparators: true,       // do not draw separators
	}}
}
