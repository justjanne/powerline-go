package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type segment struct {
	content             string
	foreground          uint8
	background          uint8
	separator           string
	separatorForeground uint8
}

type args struct {
	CwdMode          *string
	CwdMaxDepth      *int
	CwdMaxDirSize    *int
	ColorizeHostname *bool
	EastAsianWidth   *bool
	PromptOnNewLine  *bool
	Mode             *string
	Theme            *string
	Shell            *string
	Modules          *string
	PrevError        *int
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
	"cwd":    segmentCwd,
	"docker": segmentDocker,
	"exit":   segmentExitCode,
	"git":    segmentGit,
	"hg":     segmentHg,
	"host":   segmentHost,
	"jobs":   segmentJobs,
	"perms":  segmentPerms,
	"root":   segmentRoot,
	"ssh":    segmentSsh,
	"time":   segmentTime,
	"user":   segmentUser,
	"venv":   segmentVirtualEnv,
}

func main() {
	args := args{
		CwdMode: flag.String("cwd-mode", "fancy",
			"How to display the current directory\n"+
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
			"The list of modules to load. Separate with ','\n"+
				"    	(valid choices: cwd, docker, exit, git, hg, host, jobs, perms, root, ssh, time, user, venv)\n"+
				"       "),
		PrevError: flag.Int("error", 0,
			"Exit code of previously executed command"),
	}
	flag.Parse()
	powerline := NewPowerline(args, getValidCwd())

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
