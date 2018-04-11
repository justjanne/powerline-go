package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func segmentJobs(p *powerline) {
	nJobs := -1

	pppid_out, _ := exec.Command("ps", "-p", fmt.Sprintf("%d", os.Getppid()), "-oppid=").Output()
	pppid, _ := strconv.ParseInt(strings.TrimSpace(string(pppid_out)), 10, 64)
	out, _ := exec.Command("ps", "-a", "-oppid=").Output()
	processes := strings.Split(string(out), "\n")
	for _, processPpidStr := range processes {
		processPpid, _ := strconv.ParseInt(strings.TrimSpace(processPpidStr), 10, 64)
		if int(processPpid) == int(pppid) {
			nJobs++
		}
	}

	if nJobs > 0 {
		p.appendSegment("jobs", segment{
			content:    fmt.Sprintf("%d", nJobs),
			foreground: p.theme.JobsFg,
			background: p.theme.JobsBg,
		})
	}
}
