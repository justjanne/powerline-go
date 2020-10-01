package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
)

type alignment int

const (
	alignLeft alignment = iota
	alignRight
)

const (
	// MinUnsignedInteger minimum unsigned integer
	MinUnsignedInteger uint = 0
	// MaxUnsignedInteger maximum unsigned integer
	MaxUnsignedInteger = ^MinUnsignedInteger
	// MaxInteger maximum integer
	MaxInteger = int(MaxUnsignedInteger >> 1)
	// MinInteger minimum integer
	MinInteger = ^MaxInteger
)

type args struct {
	CwdMode                *string
	CwdMaxDepth            *int
	CwdMaxDirSize          *int
	ColorizeHostname       *bool
	HostnameOnlyIfSSH      *bool
	SshAlternateIcon       *bool
	EastAsianWidth         *bool
	PromptOnNewLine        *bool
	StaticPromptIndicator  *bool
	GitAssumeUnchangedSize *int64
	Mode                   *string
	Theme                  *string
	Shell                  *string
	Modules                *string
	ModulesRight           *string
	Priority               *string
	MaxWidthPercentage     *int
	TruncateSegmentWidth   *int
	PrevError              *int
	NumericExitCodes       *bool
	IgnoreRepos            *string
	ShortenGKENames        *bool
	ShortenEKSNames        *bool
	ShellVar               *string
	TrimADDomain           *bool
	PathAliases            *string
	Duration               *string
	DurationMin            *string
	Eval                   *bool
	Condensed              *bool
}

func warn(msg string) {
	print("[powerline-go]", msg)
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func getValidCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
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

var modules = map[string]func(*powerline) []pwl.Segment{
	"aws":                 segmentAWS,
	"cwd":                 segmentCwd,
	"docker":              segmentDocker,
	"docker-context":      segmentDockerContext,
	"dotenv":              segmentDotEnv,
	"duration":            segmentDuration,
	"exit":                segmentExitCode,
	"gcp":                 segmentGCP,
	"git":                 segmentGit,
	"gitlite":             segmentGitLite,
	"hg":                  segmentHg,
	"svn":                 segmentSubversion,
	"host":                segmentHost,
	"jobs":                segmentJobs,
	"kube":                segmentKube,
	"load":                segmentLoad,
	"newline":             segmentNewline,
	"perlbrew":            segmentPerlbrew,
	"plenv":               segmentPlEnv,
	"perms":               segmentPerms,
	"rbenv":               segmentRbenv,
	"root":                segmentRoot,
	"shell-var":           segmentShellVar,
	"shenv":               segmentShEnv,
	"ssh":                 segmentSSH,
	"termtitle":           segmentTermTitle,
	"terraform-workspace": segmentTerraformWorkspace,
	"time":                segmentTime,
	"node":                segmentNode,
	"user":                segmentUser,
	"venv":                segmentVirtualEnv,
	"vgo":                 segmentVirtualGo,
	"nix-shell":           segmentNixShell,
}

func comments(lines ...string) string {
	return " " + strings.Join(lines, "\n"+" ")
}

func commentsWithDefaults(lines ...string) string {
	return comments(lines...) + "\n"
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
			comments("Colorize the hostname based on a hash of itself, or use the PLGO_HOSTNAMEFG and/or PLGO_HOSTNAMEBG env vars.")),
		HostnameOnlyIfSSH: flag.Bool(
			"hostname-only-if-ssh",
			false,
			comments("Show hostname only for SSH connections")),
		SshAlternateIcon: flag.Bool(
			"alternate-ssh-icon",
			false,
			comments("Show the older, original icon for SSH connections")),
		EastAsianWidth: flag.Bool(
			"east-asian-width",
			false,
			comments("Use East Asian Ambiguous Widths")),
		PromptOnNewLine: flag.Bool(
			"newline",
			false,
			comments("Show the prompt on a new line")),
		StaticPromptIndicator: flag.Bool(
			"static-prompt-indicator",
			false,
			comments("Always show the prompt indicator with the default color, never with the error color")),
		GitAssumeUnchangedSize: flag.Int64(
			"git-assume-unchanged-size",
			2048,
			comments("Disable checking for changed/edited files in git repositories where the index is larger than this size (in KB), improves performance")),
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
				"(valid choices: aws, cwd, docker, docker-context, dotenv, duration, exit, git, gitlite, hg, host, jobs, kube, load, newline, nix-shell, node, perlbrew, perms, plenv, root, shell-var, shenv, ssh, svn, termtitle, terraform-workspace, time, user, venv, vgo)")),
		ModulesRight: flag.String(
			"modules-right",
			"",
			comments("The list of modules to load anchored to the right, for shells that support it, separated by ','",
				"(valid choices: aws, cwd, docker, docker-context, dotenv, duration, exit, git, gitlite, hg, host, jobs, kube, load, newline, nix-shell, node, perlbrew, perms, plenv, root, shell-var, shenv, ssh, svn, termtitle, terraform-workspace, time, user, venv, vgo)")),
		Priority: flag.String(
			"priority",
			"root,cwd,user,host,ssh,perms,git-branch,git-status,hg,jobs,exit,cwd-path",
			commentsWithDefaults("Segments sorted by priority, if not enough space exists, the least priorized segments are removed first. Separate with ','",
				"(valid choices: aws, cwd, docker, docker-context, dotenv, duration, exit, git, gitlite, hg, host, jobs, kube, load, newline, nix-shell, node, perlbrew, perms, plenv, root, shell-var, shenv, ssh, svn, termtitle, terraform-workspace, time, user, venv, vgo)")),
		MaxWidthPercentage: flag.Int(
			"max-width",
			0,
			comments("Maximum width of the shell that the prompt may use, in percent. Setting this to 0 disables the shrinking subsystem.")),
		TruncateSegmentWidth: flag.Int(
			"truncate-segment-width",
			16,
			commentsWithDefaults("Maximum width of a segment, segments longer than this will be shortened if space is limited. Setting this to 0 disables it.")),
		PrevError: flag.Int(
			"error",
			0,
			comments("Exit code of previously executed command")),
		NumericExitCodes: flag.Bool(
			"numeric-exit-codes",
			false,
			comments("Shows numeric exit codes for errors.")),
		IgnoreRepos: flag.String(
			"ignore-repos",
			"",
			comments("A list of git repos to ignore. Separate with ','.",
				"Repos are identified by their root directory.")),
		ShortenGKENames: flag.Bool(
			"shorten-gke-names",
			false,
			comments("Shortens names for GKE Kube clusters.")),
		ShortenEKSNames: flag.Bool(
			"shorten-eks-names",
			false,
			comments("Shortens names for EKS Kube clusters.")),
		ShellVar: flag.String(
			"shell-var",
			"",
			comments("A shell variable to add to the segments.")),
		TrimADDomain: flag.Bool(
			"trim-ad-domain",
			false,
			comments("Trim the Domainname from the AD username.")),
		PathAliases: flag.String(
			"path-aliases",
			"",
			comments("One or more aliases from a path to a short name. Separate with ','.",
				"An alias maps a path like foo/bar/baz to a short name like FBB.",
				"Specify these as key/value pairs like foo/bar/baz=FBB.",
				"Use '~' for your home dir. You may need to escape this character to avoid shell substitution.")),
		Duration: flag.String(
			"duration",
			"",
			comments("The elapsed clock-time of the previous command")),
		DurationMin: flag.String(
			"duration-min",
			"0",
			comments("The minimal time a command has to take before the duration segment is shown")),
		Eval: flag.Bool(
			"eval",
			false,
			comments("Output prompt in 'eval' format.")),
		Condensed: flag.Bool(
			"condensed",
			false,
			comments("Remove spacing between segments")),
	}
	flag.Parse()
	if strings.HasSuffix(*args.Theme, ".json") {
		jsonTheme := themes["default"]

		file, err := ioutil.ReadFile(*args.Theme)
		if err == nil {
			err = json.Unmarshal(file, &jsonTheme)
			if err == nil {
				themes[*args.Theme] = jsonTheme
			} else {
				println("Error reading theme")
				println(err.Error())
			}
		}
	}
	if strings.HasSuffix(*args.Mode, ".json") {
		modeTheme := symbolTemplates["compatible"]

		file, err := ioutil.ReadFile(*args.Mode)
		if err == nil {
			err = json.Unmarshal(file, &modeTheme)
			if err == nil {
				symbolTemplates[*args.Mode] = modeTheme
			} else {
				println("Error reading mode")
				println(err.Error())
			}
		}
	}
	priorities := map[string]int{}
	priorityList := strings.Split(*args.Priority, ",")
	for idx, priority := range priorityList {
		priorities[priority] = len(priorityList) - idx
	}

	p := newPowerline(args, getValidCwd(), priorities, alignLeft)
	if p.supportsRightModules() && p.hasRightModules() && !*args.Eval {
		panic("Flag '-modules-right' requires '-eval' mode.")
	}

	fmt.Print(p.draw())
}
