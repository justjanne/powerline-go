package main

import (
	"strings"
)

func segmentGitLite(p *powerline) {
	if len(p.ignoreRepos) > 0 {
		out, err := runGitCommand("git", "rev-parse", "--show-toplevel")
		if err != nil {
			return
		}
		out = strings.TrimSpace(out)
		if p.ignoreRepos[out] {
			return
		}
	}

	out, err := runGitCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return
	}

	status := strings.TrimSpace(out)
	var branch string

	if status != "HEAD" {
		branch = status
	} else {
		branch = getGitDetachedBranch(p)
	}

	p.appendSegment("git-branch", segment{
		content:    branch,
		foreground: p.theme.RepoCleanFg,
		background: p.theme.RepoCleanBg,
	})
}
