package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
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

func segmentKube(p *powerline) {
	home, _ := os.LookupEnv("HOME")
	defaultConfigFile := filepath.Join(home, ".kube", "config")
	kubeConfigFile := defaultConfigFile
	kubeEnv, ok := os.LookupEnv("KUBECONFIG")

	if ok {
		possibleConfigs := strings.Split(kubeEnv, ":")
		// for now just take the last one
		// TODO: find one that works
		kubeConfigFile = possibleConfigs[0]
		if len(possibleConfigs) > 1 {
			kubeConfigFile = possibleConfigs[len(possibleConfigs)-1]
		}
	}

	//	fmt.Println("Reading config file " + kubeConfigFile)
	path, err := filepath.Abs(kubeConfigFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	// parse the config
	config, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	kc := &KubeConfig{}
	err = yaml.Unmarshal(config, kc)
	if err != nil {
		fmt.Println(err)
		return
	}

	cluster := ""
	namespace := ""
	for _, c := range kc.Contexts {
		if c.Name == kc.CurrentContext {
			cluster = c.Context.Cluster
			namespace = c.Context.Namespace
			break
		}
	}

	// When you use gke your clusters may look something like gke_projectname_availability-zone_cluster-01
	// instead I want it to read as `cluster-01`
	// TODO: perhaps we should just allow some regex/substr config option for this ? or at least a toggle for gke
	//       -shorten-gke-clustersnames  or something.
	if strings.HasPrefix(cluster, "gke") {
		segments := strings.Split(cluster, "_")
		cluster = segments[len(segments)-1]
	}

	if cluster != "" {
		p.appendSegment("kube-cluster", segment{
			content:    fmt.Sprintf("⎈ %s", cluster),
			foreground: p.theme.KubeClusterFg,
			background: p.theme.KubeClusterBg,
		})
	}

	if namespace != "" {
		p.appendSegment("kube-namespace", segment{
			content:    fmt.Sprintf("%s ⎈", namespace),
			foreground: p.theme.KubeNamespaceFg,
			background: p.theme.KubeNamespaceBg,
		})
	}
}
