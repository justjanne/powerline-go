package main

import (
    "fmt"
    "strings"
	"github.com/go-git/go-git/v5"

	pwl "github.com/justjanne/powerline-go/powerline"
)

// Get the root of a repository.
func getRepoRoot(repo *git.Repository) string {
	tree, _ := repo.Worktree()
	return strings.TrimSpace(tree.Filesystem.Root())
}

func repoBranch(repo *git.Repository) string {
	ref, err := repo.Head()
	if err != nil {
		return ""
	}
	if ref.Name().IsBranch() {
		return ref.Name().Short()
	} else {
		return ref.Hash().String()[:7]
	}
}

func segmentGitLite(p *powerline) []pwl.Segment {
	repo, err := git.PlainOpenWithOptions(p.cwd, &git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: true,
	})
	if err != nil {
		return []pwl.Segment{}
	}
	if len(p.ignoreRepos) > 0 {
		root := getRepoRoot(repo)
		if p.ignoreRepos[root] {
			return []pwl.Segment{}
		}
	}

	branch := repoBranch(repo)
	return []pwl.Segment{{
		Name:       "git-branch",
        Content: fmt.Sprintf("%s %s", p.symbols.RepoBranch, branch),
		Foreground: p.theme.RepoCleanFg,
		Background: p.theme.RepoCleanBg,
	}}
}
