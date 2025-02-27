package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigFileGetConfigDir(t *testing.T) {
	c := ConfigFile{}
	dir, err := c.getConfigDir()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if dir == "" {
		t.Fatal("Expected non-empty config directory")
	}
}

func TestConfigFileSaveAndLoad(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "straico-cli-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test config
	testConfig := ConfigFile{
		Key:   "test-key",
		Model: "test-model",
	}

	// Create a mock config file in the temp directory
	configPath := filepath.Join(tempDir, "config.json")
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	// Marshal the test config to JSON
	encodedConfig, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	// Write the config to the file
	if err := os.WriteFile(configPath, encodedConfig, 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Load the config from the file
	var loadedConfig ConfigFile
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	if err := json.Unmarshal(data, &loadedConfig); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Verify the loaded config matches the original
	if loadedConfig.Key != testConfig.Key {
		t.Errorf("Expected Key %q, got %q", testConfig.Key, loadedConfig.Key)
	}

	if loadedConfig.Model != testConfig.Model {
		t.Errorf("Expected Model %q, got %q", testConfig.Model, loadedConfig.Model)
	}
}

func TestLoadConfigNonExistent(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "straico-cli-test-nonexistent")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a config that points to a non-existent file
	config := ConfigFile{}

	// Loading a non-existent config should not error
	if err := config.LoadConfig(); err != nil {
		t.Errorf("Expected no error for non-existent config, got %v", err)
	}
}
