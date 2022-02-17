package main

import (
	"fmt"
	pwl "github.com/justjanne/powerline-go/powerline"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// KubeContext holds the kubernetes context
type KubeContext struct {
	Context struct {
		Cluster   string
		Namespace string
		User      string
	}
	Name string
}

// KubeConfig is the kubernetes configuration
type KubeConfig struct {
	Contexts       []KubeContext `yaml:"contexts"`
	CurrentContext string        `yaml:"current-context"`
}

func homePath() string {
	return os.Getenv(homeEnvName())
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

func segmentKube(p *powerline) []pwl.Segment {
	paths := append(strings.Split(os.Getenv("KUBECONFIG"), ":"), path.Join(homePath(), ".kube", "config"))
	config := &KubeConfig{}
	for _, configPath := range paths {
		temp := &KubeConfig{}
		if readKubeConfig(temp, configPath) == nil {
			config.Contexts = append(config.Contexts, temp.Contexts...)
			if config.CurrentContext == "" {
				config.CurrentContext = temp.CurrentContext
			}
		}
	}

	cluster := ""
	namespace := ""
	for _, context := range config.Contexts {
		if context.Name == config.CurrentContext {
			cluster = context.Name
			namespace = context.Context.Namespace
			break
		}
	}

	// When you use gke your clusters may look something like gke_projectname_availability-zone_cluster-01
	// instead I want it to read as `cluster-01`
	// So we remove the first 3 segments of this string, if the flag is set, and there are enough segments
	if strings.HasPrefix(cluster, "gke") && p.cfg.ShortenGKENames {
		segments := strings.Split(cluster, "_")
		if len(segments) > 3 {
			cluster = strings.Join(segments[3:], "_")
		}
	}

	// When you use openshift your clusters may look something like namespace/portal-url:port/user,
	// instead I want it to read as `portal-url`.
	// So we ensure there are three segments split by / and then choose the middle part,
	// we also remove the port number from the result.
	if p.cfg.ShortenOpenshiftNames {
		segments := strings.Split(cluster, "/")
		if len(segments) == 3 {
			cluster = segments[1]
			idx := strings.IndexByte(cluster, ':')
			if idx != -1 {
				cluster = cluster[0:idx]
			}
		}
	}

	// With AWS EKS, cluster names are ARNs; it makes more sense to shorten them
	// so "eks-infra" instead of "arn:aws:eks:us-east-1:XXXXXXXXXXXX:cluster/eks-infra
	const arnRegexString string = "^arn:aws:eks:[[:alnum:]-]+:[[:digit:]]+:cluster/(.*)$"
	arnRe := regexp.MustCompile(arnRegexString)

	if arnMatches := arnRe.FindStringSubmatch(cluster); arnMatches != nil && p.cfg.ShortenEKSNames {
		cluster = arnMatches[1]
	}
	segments := []pwl.Segment{}
	// Only draw the icon once
	kubeIconHasBeenDrawnYet := false
	if cluster != "" {
		kubeIconHasBeenDrawnYet = true
		segments = append(segments, pwl.Segment{
			Name:       "kube-cluster",
			Content:    fmt.Sprintf("⎈ %s", cluster),
			Foreground: p.theme.KubeClusterFg,
			Background: p.theme.KubeClusterBg,
		})
	}

	if namespace != "" {
		content := namespace
		if !kubeIconHasBeenDrawnYet {
			content = fmt.Sprintf("⎈ %s", content)
		}
		segments = append(segments, pwl.Segment{
			Name:       "kube-namespace",
			Content:    content,
			Foreground: p.theme.KubeNamespaceFg,
			Background: p.theme.KubeNamespaceBg,
		})
	}
	return segments
}
