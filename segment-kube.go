package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func segmentKube(p *powerline) {
	var cluster string

	// TODO: detect cluster
	cluster, err := exec.Command("kubectl", "config", "current-context")
	if err != nil {
		return
	}

	if cluster == "" {
		return
	}

	tmpl := `{.contexts[?(@.name=="` + cluster + `")].context.namespace}`
	namespace, err := exec.Command("kubectl", "config", "view", "-o", "jsonpath", "--template", tmpl)

	if strings.HasPrefix(cluster, "gke") {
		segments := strings.Split(cluster, '_')
		cluster = segments[len(segments)-1]
	}

	p.appendSegment("kube", segment{
		content:    fmt.Sprintf("⎈%s/%s⎈", cluster, namespace),
		foreground: p.theme.KubeClusterFg,
		background: p.theme.KubeClusterFg,
	})
}
