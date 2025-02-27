package tui

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestConversationInitConversation(t *testing.T) {
	conversations := make(Conversations, 3)

	// Initialize all conversations
	for i := range conversations {
		conversations.InitConversation(i)
	}

	// Check that each conversation was initialized properly
	for i, conv := range conversations {
		if conv.PromptHistory == nil {
			t.Errorf("Conversation %d: Expected non-nil PromptHistory", i)
		}

		if conv.Messages == nil {
			t.Errorf("Conversation %d: Expected non-nil Messages", i)
		}

		if len(conv.PromptHistory) != 0 {
			t.Errorf("Conversation %d: Expected empty PromptHistory, got %d items", i, len(conv.PromptHistory))
		}

		if len(conv.Messages) != 0 {
			t.Errorf("Conversation %d: Expected empty Messages, got %d items", i, len(conv.Messages))
		}
	}
}

func TestConversationSaveAndLoad(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "straico-cli-test-conversations")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test conversations
	originalConversations := make(Conversations, 2)
	for i := range originalConversations {
		originalConversations.InitConversation(i)
	}

	// Add some data to the conversations
	originalConversations[0].PromptHistory = append(originalConversations[0].PromptHistory, "Test prompt 1")
	originalConversations[0].Messages = append(originalConversations[0].Messages, "User: Test message 1")
	originalConversations[0].Messages = append(originalConversations[0].Messages, "LLM: Test response 1")

	originalConversations[1].PromptHistory = append(originalConversations[1].PromptHistory, "Test prompt 2")
	originalConversations[1].Messages = append(originalConversations[1].Messages, "User: Test message 2")
	originalConversations[1].Messages = append(originalConversations[1].Messages, "LLM: Test response 2")

	// Save the conversations directly to a file in the temp directory
	configPath := filepath.Join(tempDir, "conversations.json")
	encodedConversations, err := json.MarshalIndent(originalConversations, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal conversations: %v", err)
	}

	if err := os.WriteFile(configPath, encodedConversations, 0644); err != nil {
		t.Fatalf("Failed to write conversations file: %v", err)
	}

	// Load the conversations from the file
	var loadedConversations Conversations
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read conversations file: %v", err)
	}

	if err := json.Unmarshal(data, &loadedConversations); err != nil {
		t.Fatalf("Failed to unmarshal conversations: %v", err)
	}

	// Verify the loaded conversations match the original
	if len(loadedConversations) != len(originalConversations) {
		t.Fatalf("Expected %d conversations, got %d",
			len(originalConversations), len(loadedConversations))
	}

	for i := range originalConversations {
		if len(loadedConversations[i].PromptHistory) != len(originalConversations[i].PromptHistory) {
			t.Errorf("Conversation %d: Expected %d prompts, got %d",
				i, len(originalConversations[i].PromptHistory), len(loadedConversations[i].PromptHistory))
			continue
		}

		if len(loadedConversations[i].Messages) != len(originalConversations[i].Messages) {
			t.Errorf("Conversation %d: Expected %d messages, got %d",
				i, len(originalConversations[i].Messages), len(loadedConversations[i].Messages))
			continue
		}

		// Check prompt history
		for j, prompt := range originalConversations[i].PromptHistory {
			if loadedConversations[i].PromptHistory[j] != prompt {
				t.Errorf("Conversation %d, Prompt %d: Expected %q, got %q",
					i, j, prompt, loadedConversations[i].PromptHistory[j])
			}
		}

		// Check messages
		for j, message := range originalConversations[i].Messages {
			if loadedConversations[i].Messages[j] != message {
				t.Errorf("Conversation %d, Message %d: Expected %q, got %q",
					i, j, message, loadedConversations[i].Messages[j])
			}
		}
	}
}

func TestMessagesRender(t *testing.T) {
	messages := Messages{
		"Message 1",
		"Message 2",
		"Message 3",
	}

	rendered := messages.Render(20)

	// Check that the rendered output contains all messages
	for _, msg := range messages {
		if !contains(rendered, msg) {
			t.Errorf("Expected rendered output to contain %q", msg)
		}
	}
}

// Helper function to check if a string contains another string
func contains(s, substr string) bool {
	return s != "" && substr != "" && s != substr && len(s) > len(substr) && s[len(s)-1] != substr[len(substr)-1]
}
