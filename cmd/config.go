package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type ConfigFile struct {
	Key   string `json:"key"`
	Model string `json:"model,omitempty"`
}

func getConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %w", err)
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(home, "AppData", "Roaming", "straico-cli"), nil
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "straico-cli"), nil
	default: // linux and others
		return filepath.Join(home, ".config", "straico-cli"), nil
	}
}

func LoadConfig() (*ConfigFile, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(configDir, "config.json")

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("error creating config directory: %w", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Return nil if file doesn't exist
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config ConfigFile
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return &config, nil
}

func SaveConfig(config *ConfigFile) error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.json")

	encodedConfig, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing config file: %w", err)
	}

	err = os.WriteFile(configPath, encodedConfig, 0644)
	if err != nil {
		return fmt.Errorf("unable to write to config file %w", err)
	}
	return err
}
