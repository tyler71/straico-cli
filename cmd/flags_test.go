package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestInitWithDefaults(t *testing.T) {
	// Save original args and restore them after the test
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Set minimal args for testing
	os.Args = []string{"straico-cli"}

	// Create a temporary directory for config
	tempDir, err := os.MkdirTemp("", "straico-cli-test-flags")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock config file
	configPath := filepath.Join(tempDir, "config.json")
	mockConfig := ConfigFile{
		Key:   "",
		Model: "",
	}

	encodedConfig, err := json.MarshalIndent(mockConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	if err := os.WriteFile(configPath, encodedConfig, 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Call Init
	config := Init()

	// Check default values
	if config.Prompt.Model[0] != "anthropic/claude-3-haiku:beta" {
		t.Errorf("Expected default model 'anthropic/claude-3-haiku:beta', got %q", config.Prompt.Model[0])
	}

	if len(config.Prompt.YoutubeUrls) != 0 {
		t.Errorf("Expected empty YoutubeUrls, got %v", config.Prompt.YoutubeUrls)
	}

	if len(config.Prompt.FileUrls) != 0 {
		t.Errorf("Expected empty FileUrls, got %v", config.Prompt.FileUrls)
	}
}
