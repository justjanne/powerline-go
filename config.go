package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type SymbolMap map[string]SymbolTemplate
type ShellMap map[string]ShellInfo
type ThemeMap map[string]Theme
type AliasMap map[string]string

type Config struct {
	CwdMode                string    `json:"cwd-mode"`
	CwdMaxDepth            int       `json:"cwd-max-depth"`
	CwdMaxDirSize          int       `json:"cwd-max-dir-size"`
	ColorizeHostname       bool      `json:"colorize-hostname"`
	HostnameOnlyIfSSH      bool      `json:"hostname-only-if-ssh"`
	SshAlternateIcon       bool      `json:"alternate-ssh-icon"`
	EastAsianWidth         bool      `json:"east-asian-width"`
	PromptOnNewLine        bool      `json:"newline"`
	StaticPromptIndicator  bool      `json:"static-prompt-indicator"`
	VenvNameSizeLimit      int       `json:"venv-name-size-limit"`
	Jobs                   int       `json:"-"`
	GitAssumeUnchangedSize int64     `json:"git-assume-unchanged-size"`
	GitDisableStats        []string  `json:"git-disable-stats"`
	GitMode                string    `json:"git-mode"`
	Mode                   string    `json:"mode"`
	Theme                  string    `json:"theme"`
	Shell                  string    `json:"shell"`
	Modules                []string  `json:"modules"`
	ModulesRight           []string  `json:"modules-right"`
	Priority               []string  `json:"priority"`
	MaxWidthPercentage     int       `json:"max-width-percentage"`
	TruncateSegmentWidth   int       `json:"truncate-segment-width"`
	PrevError              int       `json:"-"`
	NumericExitCodes       bool      `json:"numeric-exit-codes"`
	IgnoreRepos            []string  `json:"ignore-repos"`
	ShortenGKENames        bool      `json:"shorten-gke-names"`
	ShortenEKSNames        bool      `json:"shorten-eks-names"`
	ShortenOpenshiftNames  bool      `json:"shorten-openshift-names"`
	ShellVar               string    `json:"shell-var"`
	ShellVarNoWarnEmpty    bool      `json:"shell-var-no-warn-empty"`
	TrimADDomain           bool      `json:"trim-ad-domain"`
	PathAliases            AliasMap  `json:"path-aliases"`
	Duration               string    `json:"-"`
	DurationMin            string    `json:"duration-min"`
	DurationLowPrecision   bool      `json:"duration-low-precision"`
	Eval                   bool      `json:"eval"`
	Condensed              bool      `json:"condensed"`
	IgnoreWarnings         bool      `json:"ignore-warnings"`
	Modes                  SymbolMap `json:"modes"`
	Shells                 ShellMap  `json:"shells"`
	Themes                 ThemeMap  `json:"themes"`
	Time                   string    `json:"-"`
	ViMode                 string    `json:"vi-mode"`
}

func (mode *SymbolTemplate) UnmarshalJSON(data []byte) error {
	type Alias SymbolTemplate
	tmp := defaults.Modes[defaults.Mode]
	err := json.Unmarshal(data, (*Alias)(&tmp))
	if err == nil {
		*mode = tmp
	}
	return err
}

func (theme *Theme) UnmarshalJSON(data []byte) error {
	type Alias Theme
	tmp := defaults.Themes[defaults.Theme]
	err := json.Unmarshal(data, (*Alias)(&tmp))
	if err == nil {
		*theme = tmp
	}
	return err
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "powerline-go", "config.json")
}

func (cfg *Config) Load() error {
	path := configPath()
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil // fail silently
	}
	return json.Unmarshal(file, cfg)
}

func (cfg *Config) Save() error {
	path := configPath()
	tmp := cfg
	tmp.Themes = map[string]Theme{}
	tmp.Modes = map[string]SymbolTemplate{}
	tmp.Shells = map[string]ShellInfo{}
	data, err := json.MarshalIndent(tmp, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}
