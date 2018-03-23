package main

import (
	"os"
	"strings"
)

const ellipsis = "\u2026"

type pathSegment struct {
	path     string
	home     bool
	root     bool
	ellipsis bool
	alias    bool
}

func maybeAliasPathSegments(p *powerline, pathSegments []pathSegment) []pathSegment {
	if p.pathAliases == nil {
		return pathSegments
	}

Aliases:
	for p, alias := range p.pathAliases {
		// This turns a string like "foo/bar/baz" into an array of strings.
		path := strings.Split(p, "/")

		// If the path has 3 elements, we know we should look at pathSegments
		// in 3-element chunks.
		size := len(path)
		// If there aren't that many segments in our path we can skip to the
		// next alias.
		if size > len(pathSegments) {
			continue Aliases
		}

	Segments:
		// We want to see if that array of strings exists in pathSegments.
		for i, _ := range pathSegments {
			// This is the upper index that we would look at. So if i is 0,
			// then we'd look at pathSegments[0,1,2], then [1,2,3], etc.. If i
			// is 2, we'd look at pathSegments[2,3,4] and so on.
			max := (i + size) - 1

			// But if the upper index is out of bounds we can short-circuit
			// and move on to the next alias.
			if max > (len(pathSegments)-i)-1 {
				continue Aliases
			}

			// Then we loop over the indices in path and compare the
			// elements. If any element doesn't match we can move on to the
			// next index in pathSegments.
			for j, _ := range path {
				if path[j] != pathSegments[i+j].path {
					continue Segments
				}
			}

			// They all matched! That means we can replace this slice with our
			// alias and skip to the next alias.
			pathSegments = append(
				pathSegments[:i],
				append(
					[]pathSegment{{
						path:  alias,
						alias: true,
					}},
					pathSegments[max+1:]...,
				)...,
			)
			continue Aliases
		}
	}

	return pathSegments
}

func toString(ps []pathSegment) string {
	b := make([]string, 0)
	for _, s := range ps {
		b = append(b, s.path)
	}
	return strings.Join(b, " > ")
}

func cwdToPathSegments(p *powerline, cwd string) []pathSegment {
	pathSegments := make([]pathSegment, 0)

	home, _ := os.LookupEnv("HOME")
	if strings.HasPrefix(cwd, home) {
		pathSegments = append(pathSegments, pathSegment{
			path: "~",
			home: true,
		})
		cwd = cwd[len(home):]
	} else if cwd == "/" {
		pathSegments = append(pathSegments, pathSegment{
			path: "/",
			root: true,
		})
	}

	cwd = strings.Trim(cwd, "/")
	names := strings.Split(cwd, "/")
	if names[0] == "" {
		names = names[1:]
	}

	for _, name := range names {
		pathSegments = append(pathSegments, pathSegment{
			path: name,
		})
	}

	return maybeAliasPathSegments(p, pathSegments)
}

func maybeShortenName(p *powerline, pathSegment string) string {
	if *p.args.CwdMaxDirSize > 0 && len(pathSegment) > *p.args.CwdMaxDirSize {
		return pathSegment[:*p.args.CwdMaxDirSize]
	} else {
		return pathSegment
	}
}

func escapeVariables(p *powerline, pathSegment string) string {
	pathSegment = strings.Replace(pathSegment, `\`, p.shellInfo.escapedBackslash, -1)
	pathSegment = strings.Replace(pathSegment, "`", p.shellInfo.escapedBacktick, -1)
	pathSegment = strings.Replace(pathSegment, `$`, p.shellInfo.escapedDollar, -1)
	return pathSegment
}

func getColor(p *powerline, pathSegment pathSegment, isLastDir bool) (uint8, uint8) {
	if pathSegment.home && p.theme.HomeSpecialDisplay {
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

	if *p.args.CwdMode == "plain" {
		home, _ := os.LookupEnv("HOME")
		if strings.HasPrefix(cwd, home) {
			cwd = "~" + cwd[len(home):]
		}

		p.appendSegment("cwd", segment{
			content:    cwd,
			foreground: p.theme.CwdFg,
			background: p.theme.PathBg,
		})
	} else {
		pathSegments := cwdToPathSegments(p, cwd)

		if *p.args.CwdMode == "dironly" {
			pathSegments = pathSegments[len(pathSegments)-1:]
		} else {
			maxDepth := *p.args.CwdMaxDepth
			if maxDepth <= 0 {
				warn("Ignoring -cwd-max-depth argument since it's smaller than or equal to 0")
			} else if len(pathSegments) > maxDepth {
				var nBefore int
				if maxDepth > 2 {
					nBefore = 2
				} else {
					nBefore = maxDepth - 1
				}
				firstPart := pathSegments[:nBefore]
				secondPart := pathSegments[len(pathSegments)+nBefore-maxDepth:]

				pathSegments = make([]pathSegment, 0)
				pathSegments = append(pathSegments, firstPart...)
				pathSegments = append(pathSegments, pathSegment{
					path:     ellipsis,
					ellipsis: true,
				})
				pathSegments = append(pathSegments, secondPart...)
			}
		}

		for idx, pathSegment := range pathSegments {
			isLastDir := idx == len(pathSegments)-1
			foreground, background := getColor(p, pathSegment, isLastDir)

			segment := segment{
				content:    escapeVariables(p, maybeShortenName(p, pathSegment.path)),
				foreground: foreground,
				background: background,
			}

			if !(pathSegment.home && p.theme.HomeSpecialDisplay) && !isLastDir {
				segment.separator = p.symbolTemplates.SeparatorThin
				segment.separatorForeground = p.theme.SeparatorFg
			}

			origin := "cwd-path"
			if isLastDir {
				origin = "cwd"
			}

			p.appendSegment(origin, segment)
		}
	}
}
