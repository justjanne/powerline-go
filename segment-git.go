package main

import (
	"fmt"
	pwl "github.com/justjanne/powerline-go/powerline"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type repoStats struct {
	ahead      int
	behind     int
	untracked  int
	notStaged  int
	staged     int
	conflicted int
	stashed    int
}

func (r repoStats) dirty() bool {
	return r.untracked+r.notStaged+r.staged+r.conflicted > 0
}

func (r repoStats) any() bool {
	return r.ahead+r.behind+r.untracked+r.notStaged+r.staged+r.conflicted+r.stashed > 0
}

func addRepoStatsSegment(nChanges int, symbol string, foreground uint8, background uint8) []pwl.Segment {
	if nChanges > 0 {
		return []pwl.Segment{{
			Name:       "git-status",
			Content:    fmt.Sprintf("%d%s", nChanges, symbol),
			Foreground: foreground,
			Background: background,
		}}
	}
	return []pwl.Segment{}
}

func (r repoStats) GitSegments(p *powerline) (segments []pwl.Segment) {
	segments = append(segments, addRepoStatsSegment(r.ahead, p.symbols.RepoAhead, p.theme.GitAheadFg, p.theme.GitAheadBg)...)
	segments = append(segments, addRepoStatsSegment(r.behind, p.symbols.RepoBehind, p.theme.GitBehindFg, p.theme.GitBehindBg)...)
	segments = append(segments, addRepoStatsSegment(r.staged, p.symbols.RepoStaged, p.theme.GitStagedFg, p.theme.GitStagedBg)...)
	segments = append(segments, addRepoStatsSegment(r.notStaged, p.symbols.RepoNotStaged, p.theme.GitNotStagedFg, p.theme.GitNotStagedBg)...)
	segments = append(segments, addRepoStatsSegment(r.untracked, p.symbols.RepoUntracked, p.theme.GitUntrackedFg, p.theme.GitUntrackedBg)...)
	segments = append(segments, addRepoStatsSegment(r.conflicted, p.symbols.RepoConflicted, p.theme.GitConflictedFg, p.theme.GitConflictedBg)...)
	segments = append(segments, addRepoStatsSegment(r.stashed, p.symbols.RepoStashed, p.theme.GitStashedFg, p.theme.GitStashedBg)...)
	return
}

func addRepoStatsSymbol(nChanges int, symbol string) string {
	if nChanges > 0 {
		return symbol
	}
	return ""
}

func (r repoStats) GitSymbols(p *powerline) string {
	var info string
	info += addRepoStatsSymbol(r.ahead, p.symbols.RepoAhead)
	info += addRepoStatsSymbol(r.behind, p.symbols.RepoBehind)
	info += addRepoStatsSymbol(r.staged, p.symbols.RepoStaged)
	info += addRepoStatsSymbol(r.notStaged, p.symbols.RepoNotStaged)
	info += addRepoStatsSymbol(r.untracked, p.symbols.RepoUntracked)
	info += addRepoStatsSymbol(r.conflicted, p.symbols.RepoConflicted)
	info += addRepoStatsSymbol(r.stashed, p.symbols.RepoStashed)
	return info
}

var branchRegex = regexp.MustCompile(`^## (?P<local>\S+?)(\.{3}(?P<remote>\S+?)( \[(ahead (?P<ahead>\d+)(, )?)?(behind (?P<behind>\d+))?])?)?$`)

func groupDict(pattern *regexp.Regexp, haystack string) map[string]string {
	match := pattern.FindStringSubmatch(haystack)
	result := make(map[string]string)
	if len(match) > 0 {
		for i, name := range pattern.SubexpNames() {
			if i != 0 {
				result[name] = match[i]
			}
		}
	}
	return result
}

func gitProcessEnv() []string {
	home, _ := os.LookupEnv("HOME")
	path, _ := os.LookupEnv("PATH")
	env := map[string]string{
		"LANG": "C",
		"HOME": home,
		"PATH": path,
	}
	result := make([]string, 0)
	for key, value := range env {
		result = append(result, fmt.Sprintf("%s=%s", key, value))
	}
	return result
}

func runGitCommand(cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)
	command.Env = gitProcessEnv()
	out, err := command.Output()
	return string(out), err
}

func parseGitBranchInfo(status []string) map[string]string {
	return groupDict(branchRegex, status[0])
}

func getGitDetachedBranch(p *powerline) string {
	out, err := runGitCommand("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		out, err := runGitCommand("git", "symbolic-ref", "--short", "HEAD")
		if err != nil {
			return "Error"
		}
		return strings.SplitN(out, "\n", 2)[0]
	}
	detachedRef := strings.SplitN(out, "\n", 2)
	return fmt.Sprintf("%s %s", p.symbols.RepoDetached, detachedRef[0])
}

func parseGitStats(status []string) repoStats {
	stats := repoStats{}
	if len(status) > 1 {
		for _, line := range status[1:] {
			if len(line) > 2 {
				code := line[:2]
				switch code {
				case "??":
					stats.untracked++
				case "DD", "AU", "UD", "UA", "DU", "AA", "UU":
					stats.conflicted++
				default:
					if code[0] != ' ' {
						stats.staged++
					}

					if code[1] != ' ' {
						stats.notStaged++
					}
				}
			}
		}
	}
	return stats
}

func repoRoot(path string) (string, error) {
	out, err := runGitCommand("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func indexSize(root string) (int64, error) {
	fileInfo, err := os.Stat(path.Join(root, ".git", "index"))
	if err != nil {
		return 0, err
	}

	return fileInfo.Size(), nil
}

func segmentGit(p *powerline) []pwl.Segment {
	repoRoot, err := repoRoot(p.cwd)
	if err != nil {
		return []pwl.Segment{}
	}

	if len(p.ignoreRepos) > 0 && p.ignoreRepos[repoRoot] {
		return []pwl.Segment{}
	}

	indexSize, err := indexSize(p.cwd)
	args := []string{
		"status", "--porcelain", "-b", "--ignore-submodules",
	}

	if p.cfg.GitAssumeUnchangedSize > 0 && indexSize > (p.cfg.GitAssumeUnchangedSize*1024) {
		args = append(args, "-uno")
	}

	out, err := runGitCommand("git", args...)
	if err != nil {
		return []pwl.Segment{}
	}

	status := strings.Split(out, "\n")
	stats := parseGitStats(status)
	branchInfo := parseGitBranchInfo(status)
	var branch string

	if branchInfo["local"] != "" {
		ahead, _ := strconv.ParseInt(branchInfo["ahead"], 10, 32)
		stats.ahead = int(ahead)

		behind, _ := strconv.ParseInt(branchInfo["behind"], 10, 32)
		stats.behind = int(behind)

		branch = branchInfo["local"]
	} else {
		branch = getGitDetachedBranch(p)
	}

	if len(p.symbols.RepoBranch) > 0 {
		branch = fmt.Sprintf("%s %s", p.symbols.RepoBranch, branch)
	}

	var foreground, background uint8
	if stats.dirty() {
		foreground = p.theme.RepoDirtyFg
		background = p.theme.RepoDirtyBg
	} else {
		foreground = p.theme.RepoCleanFg
		background = p.theme.RepoCleanBg
	}

	out, err = runGitCommand("git", "rev-list", "-g", "refs/stash")
	if err == nil && len(out) > 0 {
		stats.stashed = len(strings.Split(out, "\n")) - 1
	}

	segments := []pwl.Segment{{
		Name:       "git-branch",
		Content:    branch,
		Foreground: foreground,
		Background: background,
	}}

	for _, stat := range p.cfg.GitDisableStats {
		// "ahead, behind, staged, notStaged, untracked, conflicted, stashed"
		switch stat {
		case "ahead":
			stats.ahead = 0
		case "behind":
			stats.behind = 0
		case "staged":
			stats.staged = 0
		case "notStaged":
			stats.notStaged = 0
		case "untracked":
			stats.untracked = 0
		case "conflicted":
			stats.conflicted = 0
		case "stashed":
			stats.stashed = 0
		}
	}

	if p.cfg.GitMode == "simple" {
		if stats.any() {
			segments[0].Content += " " + stats.GitSymbols(p)
		}
	} else { // fancy
		segments = append(segments, stats.GitSegments(p)...)
	}

	return segments
}
