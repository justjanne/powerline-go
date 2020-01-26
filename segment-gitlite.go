package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"strings"
)

func segmentGitLite(p *powerline) {
	ctx, cancel := newVCSContext(p)
	defer cancel()

	if len(p.ignoreRepos) > 0 {
		out, err := runGitCommand(ctx, "git", "rev-parse", "--show-toplevel")
		if err != nil {
			return
		}
		out = strings.TrimSpace(out)
		if p.ignoreRepos[out] {
			return
		}
	}

	out, err := runGitCommand(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return
	}

	status := strings.TrimSpace(out)
	var branch string

	if status != "HEAD" {
		branch = status
	} else {
		branch = getGitDetachedBranch(ctx, p)
	}

	p.appendSegment("git-branch", pwl.Segment{
		Content:    branch,
		Foreground: p.theme.RepoCleanFg,
		Background: p.theme.RepoCleanBg,
	})
}
