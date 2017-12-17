package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mattn/go-runewidth"
	"io/ioutil"
	"os"
	"strings"
)

const (
	MinUnsignedInteger uint = 0
	MaxUnsignedInteger      = ^MinUnsignedInteger
	MaxInteger         int  = int(MaxUnsignedInteger >> 1)
	MinInteger              = ^MaxInteger
)

type segment struct {
	content             string
	foreground          uint8
	background          uint8
	separator           string
	separatorForeground uint8
	priority            int
	width               int
	special             bool
}

type args struct {
	CwdMode              *string
	CwdMaxDepth          *int
	CwdMaxDirSize        *int
	ColorizeHostname     *bool
	EastAsianWidth       *bool
	PromptOnNewLine      *bool
	Mode                 *string
	Theme                *string
	Shell                *string
	Modules              *string
	Priority             *string
	MaxWidthPercentage   *int
	TruncateSegmentWidth *int
	PrevError            *int
	IgnoreRepos          *string
	ShortenGKENames      *bool
	ShellVar             *string
}

func (s segment) computeWidth() int {
	return runewidth.StringWidth(s.content) + runewidth.StringWidth(s.separator) + 2
}

func warn(msg string) {
	print("[powerline-go]", msg)
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func getValidCwd() string {
	cwd, exists := os.LookupEnv("PWD")
	if !exists {
		warn("Your current directory is invalid.")
		print("> ")
		os.Exit(1)
	}

	parts := strings.Split(cwd, string(os.PathSeparator))
	up := cwd

	for len(parts) > 0 && !pathExists(up) {
		parts = parts[:len(parts)-1]
		up = strings.Join(parts, string(os.PathSeparator))
	}
	if cwd != up {
		warn("Your current directory is invalid. Lowest valid directory: " + up)
	}
	return cwd
}

var modules = map[string](func(*powerline)){
	"aws":      segmentAWS,
	"cwd":      segmentCwd,
	"docker":   segmentDocker,
	"dotenv":   segmentDotEnv,
	"exit":     segmentExitCode,
	"git":      segmentGit,
	"gitlite":  segmentGitLite,
	"hg":       segmentHg,
	"host":     segmentHost,
	"jobs":     segmentJobs,
	"kube":     segmentKube,
	"perlbrew": segmentPerlbrew,
	"perms":    segmentPerms,
	"root":     segmentRoot,
	"ssh":      segmentSsh,
	"time":     segmentTime,
	"user":     segmentUser,
	"venv":     segmentVirtualEnv,
	"shell-var": segmentShellVar,
}

func comments(lines ...string) string {
	return " " + strings.Join(lines, "\n"+"    \t ")
}

func commentsWithDefaults(lines ...string) string {
	return comments(lines...) + "\n" + "    \t"
}

func main() {
	args := args{
		CwdMode: flag.String(
			"cwd-mode",
			"fancy",
			commentsWithDefaults("How to display the current directory",
				"(valid choices: fancy, plain, dironly)")),
		CwdMaxDepth: flag.Int(
			"cwd-max-depth",
			5,
			commentsWithDefaults("Maximum number of directories to show in path")),
		CwdMaxDirSize: flag.Int(
			"cwd-max-dir-size",
			-1,
			commentsWithDefaults("Maximum number of letters displayed for each directory in the path")),
		ColorizeHostname: flag.Bool(
			"colorize-hostname",
			false,
			comments("Colorize the hostname based on a hash of itself")),
		EastAsianWidth: flag.Bool(
			"east-asian-width",
			false,
			comments("Use East Asian Ambiguous Widths")),
		PromptOnNewLine: flag.Bool(
			"newline",
			false,
			comments("Show the prompt on a new line")),
		Mode: flag.String(
			"mode",
			"patched",
			commentsWithDefaults("The characters used to make separators between segments.",
				"(valid choices: patched, compatible, flat)")),
		Theme: flag.String(
			"theme",
			"default",
			commentsWithDefaults("Set this to the theme you want to use",
				"(valid choices: default, low-contrast)")),
		Shell: flag.String(
			"shell",
			"bash",
			commentsWithDefaults("Set this to your shell type",
				"(valid choices: bare, bash, zsh)")),
		Modules: flag.String(
			"modules",
			"venv,user,host,ssh,cwd,perms,git,hg,jobs,exit,root",
			commentsWithDefaults("The list of modules to load, separated by ','",
				"(valid choices: aws, cwd, docker, dotenv, exit, git, gitlite, hg, host, jobs, perlbrew, perms, root, ssh, time, user, venv)")),
		Priority: flag.String(
			"priority",
			"root,cwd,user,host,ssh,perms,git-branch,git-status,hg,jobs,exit,cwd-path",
			commentsWithDefaults("Segments sorted by priority, if not enough space exists, the least priorized segments are removed first. Separate with ','",
				"(valid choices: aws, cwd, cwd-path, docker, exit, git-branch, git-status, hg, host, jobs, perlbrew, perms, root, ssh, time, user, venv)")),
		MaxWidthPercentage: flag.Int(
			"max-width",
			50,
			commentsWithDefaults("Maximum width of the shell that the prompt may use, in percent. Setting this to 0 disables the shrinking subsystem.")),
		TruncateSegmentWidth: flag.Int(
			"truncate-segment-width",
			16,
			commentsWithDefaults("Minimum width of a segment, segments longer than this will be shortened if space is limited. Setting this to 0 disables it.")),
		PrevError: flag.Int(
			"error",
			0,
			comments("Exit code of previously executed command")),
		IgnoreRepos: flag.String(
			"ignore-repos",
			"",
			comments("A list of git repos to ignore. Separate with ','",
				"Repos are identified by their root directory.")),
		ShortenGKENames: flag.Bool(
			"shorten-gke-names",
			false,
			comments("Shortens names for GKE Kube clusters.")),
		ShellVar: flag.String(
			"shell-var",
			"",
			comments("A shell variable to add to the segments.")),
	}
	flag.Parse()
	if strings.HasSuffix(*args.Theme, ".json") {
		jsonTheme := themes["default"]

		file, err := ioutil.ReadFile(*args.Theme)
		if err == nil {
			json.Unmarshal(file, &jsonTheme)
		}

		themes[*args.Theme] = jsonTheme
	}
	priorities := map[string]int{}
	priorityList := strings.Split(*args.Priority, ",")
	for idx, priority := range priorityList {
		priorities[priority] = len(priorityList) - idx
	}

	powerline := NewPowerline(args, getValidCwd(), priorities)

	for _, module := range strings.Split(*powerline.args.Modules, ",") {
		elem, ok := modules[module]
		if ok {
			elem(powerline)
		} else {
			println("Module not found: " + module)
		}
	}
	fmt.Print(powerline.draw())
}
