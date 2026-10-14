package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadKubeConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kubeconfig-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config")
	content := `apiVersion: v1
kind: Config
current-context: test-context
contexts:
- name: test-context
  context:
    cluster: test-cluster
    namespace: test-ns
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test kubeconfig: %v", err)
	}

	var cfg KubeConfig
	if err := readKubeConfig(&cfg, configPath); err != nil {
		t.Fatalf("readKubeConfig failed: %v", err)
	}

	if cfg.CurrentContext != "test-context" {
		t.Errorf("expected current-context 'test-context', got '%s'", cfg.CurrentContext)
	}
	if len(cfg.Contexts) != 1 {
		t.Fatalf("expected 1 context, got %d", len(cfg.Contexts))
	}
	if cfg.Contexts[0].Name != "test-context" {
		t.Errorf("expected context name 'test-context', got '%s'", cfg.Contexts[0].Name)
	}
	if cfg.Contexts[0].Context.Cluster != "test-cluster" {
		t.Errorf("expected cluster 'test-cluster', got '%s'", cfg.Contexts[0].Context.Cluster)
	}
	if cfg.Contexts[0].Context.Namespace != "test-ns" {
		t.Errorf("expected namespace 'test-ns', got '%s'", cfg.Contexts[0].Context.Namespace)
	}
}
