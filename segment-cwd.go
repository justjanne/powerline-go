package main

import (
	"os"
	"strings"
	"fmt"
)

const ellipsis = "\u2026"

func replaceHomeDir(cwd string) string {
	home, _ := os.LookupEnv("HOME")
	if strings.HasPrefix(cwd, home) {
		return "~" + cwd[len(home):]
	} else {
		return cwd
	}
}

func splitPathIntoNames(cwd string) []string {
	names := strings.Split(cwd, string(os.PathSeparator))

	if names[0] == "" {
		names = names[1:]
	}

	if len(names) > 0 && names[0] == "" {
		return []string{"/"}
	} else {
		return names
	}
}

func requiresSpecialHomeDisplay(p *powerline, pathSegment string) bool {
	return pathSegment == "~" && p.theme.HomeSpecialDisplay
}

func maybeShortenName(p *powerline, pathSegment string) string {
	if *p.args.CwdMaxDirSize > 0 {
		return pathSegment[:*p.args.CwdMaxDirSize]
	} else {
		return pathSegment
	}
}

func getColor(p *powerline, pathSegment string, isLastDir bool) (uint8, uint8) {
	if requiresSpecialHomeDisplay(p, pathSegment) {
		return p.theme.HomeFg, p.theme.HomeBg
	} else if isLastDir {
		return p.theme.CwdFg, p.theme.PathBg
	} else {
		return p.theme.PathFg, p.theme.PathBg
	}
}

func segmentCwd(p *powerline) {
	cwd := p.cwd
	if cwd == "" {
		cwd, _ = os.LookupEnv("PWD")
	}
	cwd = replaceHomeDir(cwd)

	if *p.args.CwdMode == "plain" {
		p.appendSegment(segment{
			content:    fmt.Sprintf(" %s ", cwd),
			foreground: p.theme.CwdFg,
			background: p.theme.PathBg,
		})
	} else {
		names := splitPathIntoNames(cwd)

		if *p.args.CwdMode == "dironly" {
			names = names[len(names)-1:]
		} else {
			maxDepth := *p.args.CwdMaxDepth
			if maxDepth <= 0 {
				warn("Ignoring -cwd-max-depth argument since it's smaller than or equal to 0")
			} else if len(names) > maxDepth {
				var nBefore int
				if maxDepth > 2 {
					nBefore = 2
				} else {
					nBefore = maxDepth - 1
				}
				firstPart := names[:nBefore]
				secondPart := names[len(names)+nBefore-maxDepth:]
				names = append(append(firstPart, ellipsis), secondPart...)
			}

			for idx, pathSegment := range names {
				isLastDir := idx == len(names)-1
				foreground, background := getColor(p, pathSegment, isLastDir)

				segment := segment{
					content:    fmt.Sprintf(" %s ", maybeShortenName(p, pathSegment)),
					foreground: foreground,
					background: background,
				}

				if !requiresSpecialHomeDisplay(p, pathSegment) && !isLastDir {
					segment.separator = p.symbolTemplates.SeparatorThin
					segment.separatorForeground = p.theme.SeparatorFg
				}

				p.appendSegment(segment)
			}
		}
	}
}
