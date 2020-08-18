package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const rubyVersionFileSuffix = "/.ruby-version"
const globalVersionFileSuffix = "/.rbenv/version"

func runRbenvCommand(cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)
	out, err := command.Output()
	return string(out), err
}

// check RBENV_VERSION variable
func checkEnvForRbenvVersion() (string, error) {
	rbenvVersion := os.Getenv("RBENV_VERSION")
	if (len(rbenvVersion) > 0) {
		return rbenvVersion, nil
	} else {
		return "", errors.New("Not found in RBENV_VERSION")
	}
}

// check existence of .ruby_version in tree until root path
func checkForRubyVersionFileInTree() (string, error) {
	var (
		workingDirectory string
		err error
	)

	workingDirectory, err = os.Getwd()
	if err == nil {
		for workingDirectory != "/" {
			rubyVersion, rubyVersionErr := ioutil.ReadFile(workingDirectory + rubyVersionFileSuffix)
			if rubyVersionErr == nil {
				return strings.TrimSpace(string(rubyVersion)), nil
			}

			workingDirectory = filepath.Dir(workingDirectory)
		}
	}

	return "", errors.New("No .ruby_version file found in tree")
}

// check for global version
func checkForGlobalVersion() (string, error) {
	homeDir, _ := os.UserHomeDir()
	globalRubyVersion, err := ioutil.ReadFile(homeDir + globalVersionFileSuffix)
	if err == nil {
		return strings.TrimSpace(string(globalRubyVersion)), nil
	} else {
		return "", errors.New("No global version file found in tree")
	}
}

// retrieve rbenv version output
func checkForRbenvOutput() (string, error) {
	// spawn rbenv and print out version
	out, err := runRbenvCommand("rbenv", "version")
	if err == nil {
		items := strings.Split(out, " ")
		if len(items) > 1 {
			return items[0], nil
		}
	}

	return "", errors.New("Not found in rbenv output")
}

func segmentRbenv(p *powerline) []pwl.Segment {
	var (
		segment string
		err error
	)

	segment, err = checkEnvForRbenvVersion()
	if err != nil {
		segment, err = checkForRubyVersionFileInTree()
	}
	if err != nil {
		segment, err = checkForGlobalVersion()
	}
	if err != nil {
		segment, err = checkForRbenvOutput()
	}
	if err != nil {
		return []pwl.Segment{}
	} else {
		return []pwl.Segment{{
			Name:       "rbenv",
			Content:    segment,
			Foreground: p.theme.TimeFg,
			Background: p.theme.TimeBg,
		}}
	}
}
