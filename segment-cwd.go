package main

import (
	"os"
	"sort"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
)

const ellipsis = "\u2026"

type pathSegment struct {
	path     string
	home     bool
	root     bool
	ellipsis bool
	alias    bool
}

type byRevLength []string

func (s byRevLength) Len() int {
	return len(s)
}
func (s byRevLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byRevLength) Less(i, j int) bool {
	return len(s[i]) > len(s[j])
}

func maybeAliasPathSegments(p *powerline, pathSegments []pathSegment) []pathSegment {
	pathSeparator := string(os.PathSeparator)

	if p.cfg.PathAliases == nil || len(p.cfg.PathAliases) == 0 {
		return pathSegments
	}

	keys := make([]string, len(p.cfg.PathAliases))
	for k := range p.cfg.PathAliases {
		keys = append(keys, k)
	}
	sort.Sort(byRevLength(keys))

Aliases:
	for _, k := range keys {
		// This turns a string like "foo/bar/baz" into an array of strings.
		path := strings.Split(strings.Trim(k, pathSeparator), pathSeparator)

		// If the path has 3 elements, we know we should look at pathSegments
		// in 3-element chunks.
		size := len(path)
		// If there aren't that many segments in our path we can skip to the
		// next alias.
		if size > len(pathSegments) {
			continue Aliases
		}

		alias := p.cfg.PathAliases[k]

	Segments:
		// We want to see if that array of strings exists in pathSegments.
		for i := range pathSegments {
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
			for j := range path {
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

func cwdToPathSegments(p *powerline, cwd string) []pathSegment {
	pathSeparator := string(os.PathSeparator)
	pathSegments := make([]pathSegment, 0)

	if cwd == p.userInfo.HomeDir {
		pathSegments = append(pathSegments, pathSegment{
			path: "~",
			home: true,
		})
		cwd = ""
	} else if strings.HasPrefix(cwd, p.userInfo.HomeDir+pathSeparator) {
		pathSegments = append(pathSegments, pathSegment{
			path: "~",
			home: true,
		})
		cwd = cwd[len(p.userInfo.HomeDir):]
	} else if cwd == pathSeparator {
		pathSegments = append(pathSegments, pathSegment{
			path: pathSeparator,
			root: true,
		})
	}

	cwd = strings.Trim(cwd, pathSeparator)
	names := strings.Split(cwd, pathSeparator)
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
	if p.cfg.CwdMaxDirSize > 0 && len(pathSegment) > p.cfg.CwdMaxDirSize {
		return pathSegment[:p.cfg.CwdMaxDirSize]
	}
	return pathSegment
}

func escapeVariables(p *powerline, pathSegment string) string {
	pathSegment = strings.Replace(pathSegment, `\`, p.shell.EscapedBackslash, -1)
	pathSegment = strings.Replace(pathSegment, "`", p.shell.EscapedBacktick, -1)
	pathSegment = strings.Replace(pathSegment, `$`, p.shell.EscapedDollar, -1)
	return pathSegment
}

func getColor(p *powerline, pathSegment pathSegment, isLastDir bool) (uint8, uint8, bool) {
	if pathSegment.home && p.theme.HomeSpecialDisplay {
		return p.theme.HomeFg, p.theme.HomeBg, true
	} else if pathSegment.alias {
		return p.theme.AliasFg, p.theme.AliasBg, true
	} else if isLastDir {
		return p.theme.CwdFg, p.theme.PathBg, false
	}
	return p.theme.PathFg, p.theme.PathBg, false
}

func segmentCwd(p *powerline) (segments []pwl.Segment) {
	cwd := p.cwd

	switch p.cfg.CwdMode {
	case "plain":
		if strings.HasPrefix(cwd, p.userInfo.HomeDir) {
			cwd = "~" + cwd[len(p.userInfo.HomeDir):]
		}

		segments = append(segments, pwl.Segment{
			Name:       "cwd",
			Content:    escapeVariables(p, cwd),
			Foreground: p.theme.CwdFg,
			Background: p.theme.PathBg,
		})
	default:
		pathSegments := cwdToPathSegments(p, cwd)

		if p.cfg.CwdMode == "dironly" {
			pathSegments = pathSegments[len(pathSegments)-1:]
		} else {
			maxDepth := p.cfg.CwdMaxDepth
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

			if p.cfg.CwdMode == "semifancy" && len(pathSegments) > 1 {
				var path string
				for idx, pathSegment := range pathSegments {
					if pathSegment.home || pathSegment.alias {
						continue
					}
					path += pathSegment.path
					if idx != len(pathSegments)-1 {
						path += string(os.PathSeparator)
					}
				}
				first := pathSegments[0]
				pathSegments = make([]pathSegment, 0)
				if first.home || first.alias {
					pathSegments = append(pathSegments, first)
				}
				pathSegments = append(pathSegments, pathSegment{
					path: path,
				})
			}
		}

		for idx, pathSegment := range pathSegments {
			isLastDir := idx == len(pathSegments)-1
			foreground, background, special := getColor(p, pathSegment, isLastDir)

			segment := pwl.Segment{
				Content:    escapeVariables(p, maybeShortenName(p, pathSegment.path)),
				Foreground: foreground,
				Background: background,
			}

			if !special {
				if p.align == alignRight && p.supportsRightModules() && idx != 0 {
					segment.Separator = p.symbols.SeparatorReverseThin
					segment.SeparatorForeground = p.theme.SeparatorFg
				} else if (p.align == alignLeft || !p.supportsRightModules()) && !isLastDir {
					segment.Separator = p.symbols.SeparatorThin
					segment.SeparatorForeground = p.theme.SeparatorFg
				}
			}

			segment.Name = "cwd-path"
			if isLastDir {
				segment.Name = "cwd"
			}

			segments = append(segments, segment)
		}
	}
	return segments
}
