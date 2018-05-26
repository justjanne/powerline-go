package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
)

type gitNotFoundErr struct{}

func (gitNotFoundErr) Error() string {
	return ".git is not found"
}

func findRepoRoot() (repoPath string, head []byte, err error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", nil, err
	}
recurseUp:
	for {
		fi, err := os.Stat(filepath.Join(pwd, ".git"))
		switch {
		case err == nil:
			if fi.IsDir() {
				break
			}
			fallthrough
		case os.IsNotExist(err):
			if pwd == "/" {
				return "", nil, gitNotFoundErr{}
			}
			pwd = filepath.Join(pwd, "..")
			continue recurseUp
		default:
			return "", nil, err
		}
		head, err := ioutil.ReadFile(filepath.Join(pwd, ".git", "HEAD"))
		return pwd, head, err
	}
}

var refPrefix = []byte("ref:")

const shortHashLen = 10

func headToOutput(head []byte) string {
	if bytes.HasPrefix(head, refPrefix) {
		// It's a ref, so return branch name.
		// Example: refs: refs/heads/master
		ref := bytes.TrimSpace(head[len(refPrefix):])
		ref = ref[bytes.LastIndexByte(ref, '/')+1:]
		return string(ref)
	}

	// Treat it as a git hash. Return first 10 chars.
	hash := bytes.TrimSpace(head)
	if len(hash) < shortHashLen {
		return string(hash)
	}
	return string(hash[:shortHashLen])
}

func segmentGitLite(p *powerline) {
	repoRoot, head, err := findRepoRoot()
	if err != nil {
		return
	}
	if p.ignoreRepos[repoRoot] {
		return
	}

	p.appendSegment("git-branch", segment{
		content:    headToOutput(head),
		foreground: p.theme.RepoCleanFg,
		background: p.theme.RepoCleanBg,
	})
}
