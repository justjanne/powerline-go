package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"fmt"
	"gopkg.in/yaml.v2"
)

type KubeContext struct {
	Context struct {
		Cluster   string
		Namespace string
		User      string
	}
	Name string
}

type KubeConfig struct {
	Contexts       []KubeContext `yaml:"contexts"`
	CurrentContext string        `yaml:"current-context"`
}

func homePath() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	}
	return os.Getenv(env)
}

func readKubeConfig(config *KubeConfig, path string) (err error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return
	}
	fileContent, err := ioutil.ReadFile(absolutePath)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(fileContent, config)
	if err != nil {
		return
	}

	return
}

func segmentKube(p *powerline) {
	paths := append(strings.Split(os.Getenv("KUBECONFIG"), ":"), path.Join(homePath(), ".kube", "config"))
	config := &KubeConfig{}
	for _, configPath := range paths {
		if readKubeConfig(config, configPath) == nil {
			break
		}
	}

	cluster := ""
	namespace := ""
	for _, context := range config.Contexts {
		if context.Name == config.CurrentContext {
			cluster = context.Context.Cluster
			namespace = context.Context.Namespace
			break
		}
	}

	// When you use gke your clusters may look something like gke_projectname_availability-zone_cluster-01
	// instead I want it to read as `cluster-01`
	// So we remove the first 3 segments of this string, if the flag is set, and there are enough segments
	if strings.HasPrefix(cluster, "gke") && *p.args.ShortenGKENames {
		segments := strings.Split(cluster, "_")
		if len(segments) > 3 {
			cluster = strings.Join(segments[3:], "_")
		}
	}

	// Only draw the icon once
	kubeIconHasBeenDrawnYet := false
	if cluster != "" {
		kubeIconHasBeenDrawnYet = true
		p.appendSegment("kube-cluster", segment{
			content:    fmt.Sprintf("⎈ %s", cluster),
			foreground: p.theme.KubeClusterFg,
			background: p.theme.KubeClusterBg,
		})
	}

	if namespace != "" {
		content := namespace
		if !kubeIconHasBeenDrawnYet {
			content = fmt.Sprintf("⎈ %s", content)
		}
		p.appendSegment("kube-namespace", segment{
			content:    content,
			foreground: p.theme.KubeNamespaceFg,
			background: p.theme.KubeNamespaceBg,
		})
	}
}
