package main

import (
	"os"
	"path/filepath"
	"runtime"

	pwl "github.com/justjanne/powerline-go/powerline"
	"gopkg.in/yaml.v2"
)

type githubHost struct {
	User string `yaml:"user"`
}

func getGitHubConfigDir() (string, error) {
	if dir := os.Getenv("GH_CONFIG_DIR"); dir != "" {
		return dir, nil
	}
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, "gh"), nil
	}
	if dir := os.Getenv("AppData"); runtime.GOOS == "windows" && dir != "" {
		return filepath.Join(dir, "GitHub CLI"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "gh"), nil
}

func getGitHubUser() string {
	configDir, err := getGitHubConfigDir()
	if err != nil {
		return ""
	}

	data, err := os.ReadFile(filepath.Join(configDir, "hosts.yml"))
	if err != nil {
		return ""
	}

	var hosts map[string]githubHost
	if err := yaml.Unmarshal(data, &hosts); err != nil {
		return ""
	}

	host := os.Getenv("GH_HOST")
	if host == "" {
		host = "github.com"
	}

	return hosts[host].User
}

func segmentGitHub(p *powerline) []pwl.Segment {
	user := getGitHubUser()
	if user == "" {
		return []pwl.Segment{}
	}

	return []pwl.Segment{{
		Name:       "github",
		Content:    user,
		Foreground: p.theme.GitHubFg,
		Background: p.theme.GitHubBg,
	}}
}
