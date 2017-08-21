package main

import (
	"fmt"
	"os"
	"strings"
	"flag"
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
	Mode             *string
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
	"cwd":   segmentCwd,
	"exit":  segmentExitCode,
	"git":   segmentGit,
	"hg":    segmentHg,
	"host":  segmentHost,
	"jobs":  segmentJobs,
	"perms": segmentPerms,
	"root":  segmentRoot,
	"ssh":   segmentSsh,
	"user":  segmentUser,
	"venv":  segmentVirtualEnv,
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
		Mode: flag.String("mode", "patched",
			"The characters used to make separators between segments.\n"+
				"    	(valid choices: patched, compatible, flat)\n"+
				"       "),
		Shell: flag.String("shell", "bash",
			"Set this to your shell type\n"+
				"    	(valid choices: bare, bash, zsh)\n"+
				"       "),
		Modules: flag.String("modules",
			"venv,user,ssh,host,cwd,perms,jobs,exit,root",
			"The list of modules to load. Separate with ','\n"+
				"    	(valid choices: cwd, exit, git, hg, host, jobs, perms, root, ssh, user, venv)\n"+
				"       "),
		PrevError: flag.Int("error", 0,
			"Exit code of previously executed command"),
	}
	flag.Parse()
	powerline := NewPowerline(args, getValidCwd(), defaultTheme)

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
