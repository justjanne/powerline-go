package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
	"gopkg.in/yaml.v2"
)

const ppFile = "Pulumi.yaml"

type Project struct {
	Name        string
	Runtime     string
	Description string
}

type Workspace struct {
	Stack string
}

func segmentPulumiStack(p *powerline) []pwl.Segment {
	// Check for presence of Pulumi.yaml in current directory
	projDir, err := os.Getwd()
	projFile := filepath.Join(projDir, ppFile)
	stat, err := os.Stat(projFile)
	if err != nil {
		return []pwl.Segment{}
	}
	if stat.IsDir() {
		return []pwl.Segment{}
	}

	// Get active stack from project workspace JSON file
	pulumiHome := os.Getenv("PULUMI_HOME")
	if pulumiHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return []pwl.Segment{}
		}
		pulumiHome = homeDir + "/.pulumi"
	}
	y, err := ioutil.ReadFile(projFile)
	if err != nil {
		return []pwl.Segment{}
	}
	var projectInfo Project
	err = yaml.Unmarshal(y, &projectInfo)
	projName := projectInfo.Name
	hash := sha1.New()
	hash.Write([]byte(projFile))
	sha1_hash := hex.EncodeToString(hash.Sum(nil))
	wsFilename := pulumiHome + "/workspaces/" + projName + "-" + sha1_hash + "-workspace.json"
	j, err := ioutil.ReadFile(wsFilename)
	if err != nil {
		return []pwl.Segment{}
	}
	var currentWs Workspace
	err = json.Unmarshal(j, &currentWs)
	activeStack := currentWs.Stack

	return []pwl.Segment{{
		Name:       "pulumi-stack",
		Content:    strings.Replace(string(activeStack), "\n", "", -1),
		Foreground: p.theme.TFWsFg,
		Background: p.theme.TFWsBg,
	}}
}
