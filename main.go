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
	IgnoreRepos          *string
	PrevError            *int
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
	"perlbrew": segmentPerlbrew,
	"perms":    segmentPerms,
	"root":     segmentRoot,
	"ssh":      segmentSsh,
	"time":     segmentTime,
	"user":     segmentUser,
	"venv":     segmentVirtualEnv,
	"kube":     segmentKube,
}

func main() {
	args := args{
		CwdMode: flag.String("cwd-mode", "fancy",
			"How to display the current directory\n"+
				"    	(valid choices: fancy, plain, dironly)\n"+
				"       "),
		CwdMaxDepth: flag.Int("cwd-max-depth", 5,
			"Maximum number of directories to show in path\n"+
				"       "),
		CwdMaxDirSize: flag.Int("cwd-max-dir-size", -1,
			"Maximum number of letters displayed for each directory in the path\n"+
				"       "),
		ColorizeHostname: flag.Bool("colorize-hostname", false,
			"Colorize the hostname based on a hash of itself"),
		EastAsianWidth: flag.Bool("east-asian-width", false,
			"Use East Asian Ambiguous Widths"),
		PromptOnNewLine: flag.Bool("newline", false,
			"Show the prompt on a new line"),
		Mode: flag.String("mode", "patched",
			"The characters used to make separators between segments.\n"+
				"    	(valid choices: patched, compatible, flat)\n"+
				"       "),
		Theme: flag.String("theme", "default",
			"Set this to the theme you want to use\n"+
				"    	(valid choices: default, low-contrast)\n"+
				"       "),
		Shell: flag.String("shell", "bash",
			"Set this to your shell type\n"+
				"    	(valid choices: bare, bash, zsh)\n"+
				"       "),
		Modules: flag.String("modules",
			"venv,user,host,ssh,cwd,perms,git,hg,jobs,exit,root",
			"The list of modules to load, separated by ','\n"+
				"    	(valid choices: aws, cwd, docker, dotenv, exit, git, gitlite, hg, host, jobs, kube, perlbrew, perms, root, ssh, time, user, venv)\n"+
				"       "),
		Priority: flag.String("priority",
			"root,cwd,user,host,ssh,perms,git-branch,git-status,hg,jobs,exit",
			"Segments sorted by priority, if not enough space exists, the least priorized segments are removed first. Separate with ','\n"+
				"    	(valid choices: aws, cwd, docker, exit, git-branch, git-status, hg, host, jobs, kube, perlbrew, perms, root, ssh, time, user, venv)\n"+
				"       "),
		MaxWidthPercentage: flag.Int("max-width",
			50,
			"Maximum width of the shell that the prompt may use, in percent. Setting this to 0 disables the shrinking subsystem.\n"+
				"       "),
		TruncateSegmentWidth: flag.Int("truncate-segment-width",
			16,
			"Minimum width of a segment, segments longer than this will be shortened if space is limited. Setting this to 0 disables it.\n"+
				"       "),
		PrevError: flag.Int("error", 0,
			"Exit code of previously executed command"),
		IgnoreRepos: flag.String("ignore-repos",
			"",
			"A list of git repos to ignore. Separate with ','\n"+
				"Repos are identified by their root directory."),
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
