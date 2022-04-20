package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
)

const gcloudCoreSectionHeader = "\n[core]\n"

func getCloudConfigDir() (string, error) {
	p, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if runtime.GOOS != "windows" {
		p += "/.config"
	}
	p += "/gcloud"
	return p, nil
}

func getActiveGCloudConfig(configDir string) (string, error) {
	activeConfigPath := configDir + "/active_config"

	stat, err := os.Stat(activeConfigPath)
	if (err == nil && os.IsNotExist(err)) || (err == nil && stat.IsDir()) {
		return "default", nil
	} else if err != nil {
		return "", err
	}

	contents, err := ioutil.ReadFile(activeConfigPath)
	if err != nil {
		return "", err
	}

	config := strings.TrimSpace(string(contents))
	if config == "" {
		config = "default"
	}

	return config, nil
}

func getGCPProjectFromGCloud() (string, error) {
	out, err := exec.Command("gcloud", "config", "list", "project", "--format", "value(core.project)").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(string(out), "\n"), nil
}

func getGCPProjectFromFile() (string, error) {
	configDir, err := getCloudConfigDir()
	if err != nil {
		return "", err
	}

	activeConfig, err := getActiveGCloudConfig(configDir)
	if err != nil {
		return "", err
	}

	configPath := configDir + "/configurations/config_" + activeConfig
	stat, err := os.Stat(configPath)
	if err != nil {
		return "", err
	} else if stat.IsDir() {
		return "", fmt.Errorf("%s is a directory", configPath)
	}

	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return "", err
	}
	b = append([]byte("\n"), b...)

	coreStart := bytes.Index(b, []byte(gcloudCoreSectionHeader))
	if coreStart == -1 {
		return "", fmt.Errorf("could not find [core] section in %s", configPath)
	}
	b = b[coreStart+len(gcloudCoreSectionHeader):]

	coreEnd := bytes.Index(b, []byte("\n["))
	if coreEnd != -1 {
		b = b[:coreEnd]
	}

	lines := bytes.Split(b[coreStart+len(gcloudCoreSectionHeader):coreEnd], []byte("\n"))
	for _, line := range lines {
		parts := bytes.Split(line, []byte("="))
		if len(parts) == 2 {
			if strings.TrimSpace(string(parts[0])) == "project" {
				return strings.TrimSpace(string(parts[1])), nil
			}
		}
	}

	return "", nil
}

func getGCPProject() (string, error) {
	if project, err := getGCPProjectFromFile(); err == nil {
		return project, nil
	} else {
		return getGCPProjectFromGCloud()
	}
}

func segmentGCP(p *powerline) []pwl.Segment {
	project, err := getGCPProject()
	if err != nil {
		log.Fatal(err)
	}

	if project == "" {
		return []pwl.Segment{}
	}
	return []pwl.Segment{{
		Name:       "gcp",
		Content:    project,
		Foreground: p.theme.GCPFg,
		Background: p.theme.GCPBg,
	}}
}
