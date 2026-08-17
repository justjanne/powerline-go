package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
)

// utf8BOM is the UTF-8 byte order mark that some editors may add
var natsUtf8BOM = []byte{0xEF, 0xBB, 0xBF}

type natsContextInfo struct {
	Name string `json:"name"`
}

func getNatsConfigDir() string {
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "nats")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	if runtime.GOOS == "windows" {
		return filepath.Join(home, "AppData", "Local", "nats")
	}

	return filepath.Join(home, ".config", "nats")
}

func getNatsContextName() string {
	if ctx := os.Getenv("NATS_CONTEXT"); ctx != "" {
		return strings.TrimSpace(ctx)
	}

	configDir := getNatsConfigDir()
	if configDir != "" {
		contextFile := filepath.Join(configDir, "context.txt")
		if data, err := os.ReadFile(contextFile); err == nil {
			data = bytes.TrimPrefix(data, natsUtf8BOM)
			if name := strings.TrimSpace(string(data)); name != "" {
				return name
			}
		}
	}

	out, err := exec.Command("nats", "context", "info", "--json").Output()
	if err != nil {
		return ""
	}

	var ctx natsContextInfo
	if err := json.Unmarshal(out, &ctx); err != nil {
		return ""
	}

	return strings.TrimSpace(ctx.Name)
}

func segmentNats(p *powerline) []pwl.Segment {
	name := getNatsContextName()
	if name == "" {
		return []pwl.Segment{}
	}

	return []pwl.Segment{{
		Name:       "nats",
		Content:    name,
		Foreground: p.theme.NatsFg,
		Background: p.theme.NatsBg,
	}}
}
