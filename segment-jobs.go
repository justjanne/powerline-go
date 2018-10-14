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

	ppid := os.Getppid()
	if *p.args.Shell == "bash" {
		pppidOut, _ := exec.Command("ps", "-p", strconv.Itoa(ppid), "-oppid=").Output()
		pppid, _ := strconv.ParseInt(strings.TrimSpace(string(pppidOut)), 10, 64)
		ppid = int(pppid)
	}

	out, _ := exec.Command("ps", "-a", "-oppid=").Output()
	processes := strings.Split(string(out), "\n")
	for _, processPpidStr := range processes {
		processPpid, _ := strconv.ParseInt(strings.TrimSpace(processPpidStr), 10, 64)
		if int(processPpid) == ppid {
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
