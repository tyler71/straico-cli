package prompt

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPromptRead(t *testing.T) {
	p := Prompt{
		Message: "Test message",
		Model:   []string{"test-model"},
	}

	// Create a buffer to read into
	buffer := make([]byte, 1024)

	n, err := p.Read(buffer)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if n == 0 {
		t.Fatal("Expected to read some bytes, got 0")
	}

	// Unmarshal the buffer back into a Prompt to verify
	var readPrompt Prompt
	if err := json.Unmarshal(buffer[:n], &readPrompt); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if readPrompt.Message != p.Message {
		t.Errorf("Expected Message %q, got %q", p.Message, readPrompt.Message)
	}

	if len(readPrompt.Model) != len(p.Model) || readPrompt.Model[0] != p.Model[0] {
		t.Errorf("Expected Model %v, got %v", p.Model, readPrompt.Model)
	}
}

func TestPromptRequest(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check authorization header
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("Expected Authorization header 'Bearer test-key', got %q", r.Header.Get("Authorization"))
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		mockResponse := `{
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
		}`
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Create a custom request function that uses the test server
	requestWithURL := func(p Prompt, key, text, url string, context []string) (StraicoResponse, error) {
		promptHistory := strings.Join(context, "\n")
		contextLength := len(promptHistory)
		if contextLength < MaxContextLength {
			contextLength = MaxContextLength
		}
		if len(context) > 1 {
			if contextLength > 1000 {
				p.Message = "Answer the question using the context below.\n" + promptHistory[contextLength-1000:] + "\nQuestion:" + text + "\nAnswer:"
			} else {
				p.Message = "Answer the question using the context below.\n" + promptHistory + "\nQuestion:" + text + "\nAnswer:"
			}
		} else {
			p.Message = text
		}
		jsonAbc, _ := json.Marshal(p)
		client := &http.Client{}
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonAbc))

		req.Header = http.Header{
			"Authorization": []string{"Bearer " + key},
			"Content-Type":  []string{"application/json"},
			"Accept":        []string{"application/json"},
		}
		resp, _ := client.Do(req)
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			return StraicoResponse{}, nil
		}
		bodyText, err := io.ReadAll(resp.Body)
		if err != nil {
			return StraicoResponse{}, err
		}

		return UnmarshalStraicoResponse(bodyText)
	}

	// Create a prompt and make a request using the test server
	p := Prompt{
		Message: "Test message",
		Model:   []string{"test-model"},
	}

	response, err := requestWithURL(p, "test-key", "Test message", server.URL, []string{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !response.Success {
		t.Error("Expected Success to be true")
	}

	if response.Data.OverallPrice.Total != 0.3 {
		t.Errorf("Expected Total price 0.3, got %f", response.Data.OverallPrice.Total)
	}

	completion := response.Data.Completions["test-model"].Completion
	if len(completion.Choices) == 0 {
		t.Fatal("Expected at least one choice")
	}

	content := completion.Choices[0].Message.Content
	if content != "Test response" {
		t.Errorf("Expected content 'Test response', got %q", content)
	}
}
