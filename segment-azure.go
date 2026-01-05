package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
)

// utf8BOM is the UTF-8 byte order mark that Azure CLI may add to JSON files
var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

type azureSubscription struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	State     string `json:"state"`
	IsDefault bool   `json:"isDefault"`
}

type azureProfile struct {
	Subscriptions []azureSubscription `json:"subscriptions"`
}

func getAzureSubscription() string {
	envSubID := os.Getenv("AZURE_SUBSCRIPTION_ID")

	home, err := os.UserHomeDir()
	if err != nil {
		return envSubID
	}

	data, err := os.ReadFile(filepath.Join(home, ".azure", "azureProfile.json"))
	if err != nil {
		return envSubID
	}

	data = bytes.TrimPrefix(data, utf8BOM)

	var profile azureProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return envSubID
	}

	if envSubID != "" {
		for _, sub := range profile.Subscriptions {
			if sub.ID == envSubID {
				return sub.Name
			}
		}
		return envSubID
	}

	var firstEnabled string
	for _, sub := range profile.Subscriptions {
		if sub.State != "Enabled" {
			continue
		}
		if sub.IsDefault {
			return sub.Name
		}
		if firstEnabled == "" {
			firstEnabled = sub.Name
		}
	}

	return firstEnabled
}

func getAzureResourceGroup() string {
	if rg := os.Getenv("AZURE_DEFAULTS_GROUP"); rg != "" {
		return rg
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	data, err := os.ReadFile(filepath.Join(home, ".azure", "config"))
	if err != nil {
		return ""
	}

	inDefaults := false
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)

		if line == "[defaults]" {
			inDefaults = true
			continue
		}

		if strings.HasPrefix(line, "[") {
			inDefaults = false
			continue
		}

		if inDefaults && (strings.HasPrefix(line, "group ") || strings.HasPrefix(line, "group=")) {
			if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}

	return ""
}

func segmentAzure(p *powerline) []pwl.Segment {
	subscription := getAzureSubscription()
	if subscription == "" {
		return []pwl.Segment{}
	}

	content := subscription
	if rg := getAzureResourceGroup(); rg != "" {
		content += " (" + rg + ")"
	}

	return []pwl.Segment{{
		Name:       "azure",
		Content:    content,
		Foreground: p.theme.AzureFg,
		Background: p.theme.AzureBg,
	}}
}
