package tui

import (
	"encoding/json"
	"fmt"
	"github.com/tyler71/straico-cli/m/v0/prompt"
	"os"
	"path/filepath"
	"runtime"
)

const saveFile = "conversations.json"

type Conversation struct {
	PromptHistory []string `json:"prompt_history"`
	Messages      Messages `json:"messages"`
}
type Conversations []Conversation

func (c Conversations) InitConversation(channel int) {
	c[channel] = Conversation{
		PromptHistory: make([]string, 0, prompt.MaxContextLength),
		Messages:      make(Messages, 0, prompt.MaxContextLength*2),
	}
}

func (c Conversations) getConfigDir() (string, error) {
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

func (c Conversations) LoadConversations() error {
	configDir, err := c.getConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, saveFile)

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Return nil if file doesn't exist
		}
		return fmt.Errorf("error reading config file: %w", err)
	}

	if err := json.Unmarshal(data, &c); err != nil {
		return fmt.Errorf("error parsing config file: %w", err)
	}

	return nil
}

func (c *Conversations) SaveConversations() error {
	configDir, err := c.getConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, saveFile)

	encodedConfig, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing config file: %w", err)
	}

	err = os.WriteFile(configPath, encodedConfig, 0644)
	if err != nil {
		return fmt.Errorf("unable to write to config file %w", err)
	}
	return err
}
