package prompt

import (
	"testing"
)

func TestUnmarshalStraicoResponse(t *testing.T) {
	jsonData := []byte(`{
		"data": {
			"overall_price": {
				"input": 0.1,
				"output": 0.2,
				"total": 0.3
			},
			"overall_words": {
				"input": 10,
				"output": 20,
				"total": 30
			},
			"completions": {
				"test-model": {
					"completion": {
						"id": "test-id",
						"model": "test-model",
						"object": "test-object",
						"created": 1234567890,
						"choices": [
							{
								"index": 0,
								"message": {
									"role": "assistant",
									"content": "Test response"
								},
								"finish_reason": "stop",
								"logprobs": null
							}
						],
						"usage": {
							"prompt_tokens": 10,
							"completion_tokens": 20,
							"total_tokens": 30
						}
					},
					"price": {
						"input": 0.1,
						"output": 0.2,
						"total": 0.3
					},
					"words": {
						"input": 10,
						"output": 20,
						"total": 30
					}
				}
			}
		},
		"success": true
	}`)

	response, err := UnmarshalStraicoResponse(jsonData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !response.Success {
		t.Error("Expected Success to be true")
	}

	if response.Data.OverallPrice.Total != 0.3 {
		t.Errorf("Expected OverallPrice.Total 0.3, got %f", response.Data.OverallPrice.Total)
	}

	if response.Data.OverallWords.Total != 30 {
		t.Errorf("Expected OverallWords.Total 30, got %f", response.Data.OverallWords.Total)
	}

	completion, exists := response.Data.Completions["test-model"]
	if !exists {
		t.Fatal("Expected completion for 'test-model' to exist")
	}

	if completion.Completion.Model != "test-model" {
		t.Errorf("Expected Completion.Model 'test-model', got %q", completion.Completion.Model)
	}

	if len(completion.Completion.Choices) != 1 {
		t.Fatalf("Expected 1 choice, got %d", len(completion.Completion.Choices))
	}

	choice := completion.Completion.Choices[0]
	if choice.Message.Role != "assistant" {
		t.Errorf("Expected Message.Role 'assistant', got %q", choice.Message.Role)
	}

	if choice.Message.Content != "Test response" {
		t.Errorf("Expected Message.Content 'Test response', got %q", choice.Message.Content)
	}
}
