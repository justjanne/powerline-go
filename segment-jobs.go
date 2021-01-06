package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
	"os/exec"
	"strings"
	"strconv"
)

func segmentJobs(p *powerline) []pwl.Segment {
	ppid := strconv.Itoa(os.Getppid())
	out, _ := exec.Command("ps", "--ppid", ppid, "-opid=").Output()
	nJobs := strings.Count(string(out), "\n") - 1

	if nJobs > 0 {
		return []pwl.Segment{{
			Name:       "jobs",
			Content:    strconv.Itoa(nJobs),
			Foreground: p.theme.JobsFg,
			Background: p.theme.JobsBg,
		}}
	}
	return []pwl.Segment{}
}
