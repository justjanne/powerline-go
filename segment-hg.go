package main

import (
	"context"
	"fmt"
	pwl "github.com/justjanne/powerline-go/powerline"
	"os/exec"
	"strings"
)

func getHgStatus(ctx context.Context) (bool, bool, bool) {
	hasModifiedFiles := false
	hasUntrackedFiles := false
	hasMissingFiles := false

	out, err := exec.CommandContext(ctx, "hg", "status").Output()
	if err == nil {
		output := strings.Split(string(out), "\n")
		for _, line := range output {
			if line != "" {
				if line[0] == '?' {
					hasUntrackedFiles = true
				} else if line[0] == '!' {
					hasMissingFiles = true
				} else {
					hasModifiedFiles = true
				}
			}
		}
	}
	return hasModifiedFiles, hasUntrackedFiles, hasMissingFiles
}

func segmentHg(p *powerline) {
	ctx, cancel := newVCSContext(p)
	defer cancel()

	out, _ := exec.CommandContext(ctx, "hg", "branch").Output()
	output := strings.SplitN(string(out), "\n", 2)
	if len(output) > 0 && output[0] != "" {
		branch := output[0]
		hasModifiedFiles, hasUntrackedFiles, hasMissingFiles := getHgStatus(ctx)

		var foreground, background uint8
		var content string
		if hasModifiedFiles || hasUntrackedFiles || hasMissingFiles {
			foreground = p.theme.RepoDirtyFg
			background = p.theme.RepoDirtyBg

			extra := ""

			if hasUntrackedFiles {
				extra += "+"
			}

			if hasMissingFiles {
				extra += "!"
			}

			content = fmt.Sprintf("%s %s", branch, extra)
		} else {
			foreground = p.theme.RepoCleanFg
			background = p.theme.RepoCleanBg

			content = fmt.Sprintf(branch)
		}

		p.appendSegment("hg", pwl.Segment{
			Content:    content,
			Foreground: foreground,
			Background: background,
		})
	}
}
