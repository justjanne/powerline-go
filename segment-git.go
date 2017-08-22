package main

import (
	"fmt"
	"os"
	"os/exec"
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
}

func (r repoStats) dirty() bool {
	return r.untracked+r.notStaged+r.staged+r.conflicted > 0
}

func addRepoStatsSegment(p *powerline, nChanges int, symbol string, foreground uint8, background uint8) {
	if nChanges > 0 {
		p.appendSegment(segment{
			content:    fmt.Sprintf(" %d%s ", nChanges, symbol),
			foreground: foreground,
			background: background,
		})
	}
}

func (r repoStats) addToPowerline(p *powerline) {
	addRepoStatsSegment(p, r.ahead, p.symbolTemplates.RepoAhead, p.theme.GitAheadFg, p.theme.GitAheadBg)
	addRepoStatsSegment(p, r.behind, p.symbolTemplates.RepoBehind, p.theme.GitBehindFg, p.theme.GitBehindBg)
	addRepoStatsSegment(p, r.staged, p.symbolTemplates.RepoStaged, p.theme.GitStagedFg, p.theme.GitStagedBg)
	addRepoStatsSegment(p, r.notStaged, p.symbolTemplates.RepoNotStaged, p.theme.GitNotStagedFg, p.theme.GitNotStagedBg)
	addRepoStatsSegment(p, r.untracked, p.symbolTemplates.RepoUntracked, p.theme.GitUntrackedFg, p.theme.GitUntrackedBg)
	addRepoStatsSegment(p, r.conflicted, p.symbolTemplates.RepoConflicted, p.theme.GitConflictedFg, p.theme.GitConflictedBg)
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

func parseGitBranchInfo(status []string) map[string]string {
	return groupDict(branchRegex, status[0])
}

func getGitDetachedBranch(p *powerline) string {
	command := exec.Command("git", "describe", "--tags", "--always")
	command.Env = gitProcessEnv()
	out, err := command.Output()
	if err != nil {
		return "Error"
	} else {
		detachedRef := strings.SplitN(string(out), "\n", 2)
		return fmt.Sprintf("%s %s", p.symbolTemplates.RepoDetached, detachedRef)
	}
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

func segmentGit(p *powerline) {
	command := exec.Command("git", "status", "--porcelain", "-b")
	command.Env = gitProcessEnv()
	out, err := command.Output()
	if err != nil {
	} else {
		status := strings.Split(string(out), "\n")
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

		var foreground, background uint8
		if stats.dirty() {
			foreground = p.theme.RepoDirtyFg
			background = p.theme.RepoDirtyBg
		} else {
			foreground = p.theme.RepoCleanFg
			background = p.theme.RepoCleanBg
		}

		p.appendSegment(segment{
			content:    fmt.Sprintf(" %s ", branch),
			foreground: foreground,
			background: background,
		})
		stats.addToPowerline(p)
	}
}
