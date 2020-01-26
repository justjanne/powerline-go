package main

import (
	"context"
	"errors"
	"fmt"
	pwl "github.com/justjanne/powerline-go/powerline"
	"os/exec"
	"strings"
)

var otherModified int

func addSvnRepoStatsSegment(p *powerline, nChanges int, symbol string, foreground uint8, background uint8) {
	if nChanges > 0 {
		p.appendSegment("svn-status", pwl.Segment{
			Content:    fmt.Sprintf("%d%s", nChanges, symbol),
			Foreground: foreground,
			Background: background,
		})
	}
}

func (r repoStats) addSvnToPowerline(p *powerline) {
	addSvnRepoStatsSegment(p, r.ahead, p.symbolTemplates.RepoAhead, p.theme.GitAheadFg, p.theme.GitAheadBg)
	addSvnRepoStatsSegment(p, r.behind, p.symbolTemplates.RepoBehind, p.theme.GitBehindFg, p.theme.GitBehindBg)
	addSvnRepoStatsSegment(p, r.staged, p.symbolTemplates.RepoStaged, p.theme.GitStagedFg, p.theme.GitStagedBg)
	addSvnRepoStatsSegment(p, r.notStaged, p.symbolTemplates.RepoNotStaged, p.theme.GitNotStagedFg, p.theme.GitNotStagedBg)
	addSvnRepoStatsSegment(p, r.untracked, p.symbolTemplates.RepoUntracked, p.theme.GitUntrackedFg, p.theme.GitUntrackedBg)
	addSvnRepoStatsSegment(p, r.conflicted, p.symbolTemplates.RepoConflicted, p.theme.GitConflictedFg, p.theme.GitConflictedBg)
	addSvnRepoStatsSegment(p, r.stashed, p.symbolTemplates.RepoStashed, p.theme.GitStashedFg, p.theme.GitStashedBg)
}

func runSvnCommand(ctx context.Context, cmd string, args ...string) (string, error) {
	command := exec.CommandContext(ctx, cmd, args...)
	out, err := command.Output()
	return string(out), err
}

func parseSvnURL(ctx context.Context) (map[string]string, error) {
	info, err := runSvnCommand(ctx, "svn", "info")
	if err != nil {
		return nil, errors.New("not a working copy")
	}

	svnInfo := make(map[string]string, 0)
	infos := strings.Split(info, "\n")
	if len(infos) > 1 {
		for _, line := range infos[:] {
			items := strings.Split(line, ": ")
			if len(items) >= 2 {
				svnInfo[items[0]] = items[1]
			}
		}
	}

	return svnInfo, nil
}

func ensureUnmodified(code string, stats repoStats) {
	if code != " " {
		otherModified++
	}
}

func parseSvnStatus(ctx context.Context) repoStats {
	stats := repoStats{}
	info, err := runSvnCommand(ctx, "svn", "status", "-u")
	if err != nil {
		return stats
	}
	infos := strings.Split(info, "\n")
	if len(infos) > 1 {
		for _, line := range infos[:] {
			if len(line) >= 9 {
				code := line[0:1]
				switch code {
				case "?":
					stats.untracked++
				case "C":
					stats.conflicted++
				case "A", "D", "M":
					stats.notStaged++
				default:
					ensureUnmodified(code, stats)
				}
				code = line[1:2]
				switch code {
				case "C":
					stats.conflicted++
				case "M":
					stats.notStaged++
				default:
					ensureUnmodified(code, stats)
				}
				ensureUnmodified(line[2:3], stats)
				ensureUnmodified(line[3:4], stats)
				ensureUnmodified(line[4:5], stats)
				ensureUnmodified(line[5:6], stats)
				ensureUnmodified(line[6:7], stats)
				ensureUnmodified(line[7:8], stats)
				code = line[8:9]
				switch code {
				case "*":
					stats.behind++
				default:
					ensureUnmodified(code, stats)
				}
			}
		}
	}

	return stats
}

func segmentSubversion(p *powerline) {
	ctx, cancel := newVCSContext(p)
	defer cancel()

	svnInfo, err := parseSvnURL(ctx)
	if err != nil {
		return
	}

	if len(p.ignoreRepos) > 0 {
		if p.ignoreRepos[svnInfo["URL"]] || p.ignoreRepos[svnInfo["Relative URL"]] {
			return
		}
	}

	svnStats := parseSvnStatus(ctx)

	var foreground, background uint8
	if svnStats.dirty() || otherModified > 0 {
		foreground = p.theme.RepoDirtyFg
		background = p.theme.RepoDirtyBg
	} else {
		foreground = p.theme.RepoCleanFg
		background = p.theme.RepoCleanBg
	}

	p.appendSegment("svn-branch", pwl.Segment{
		Content:    svnInfo["Relative URL"],
		Foreground: foreground,
		Background: background,
	})

	svnStats.addSvnToPowerline(p)
}
