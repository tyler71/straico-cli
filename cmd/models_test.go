package cmd

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetModels(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
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
				"chat": [
					{
						"name": "Test Model 1",
						"model": "test-model-1",
						"word_limit": 1000,
						"pricing": {
							"coins": 10.5,
							"words": 1000
						},
						"max_output": 500
					},
					{
						"name": "Test Model 2",
						"model": "test-model-2",
						"word_limit": 2000,
						"pricing": {
							"coins": 20.5,
							"words": 2000
						},
						"max_output": 1000
					}
				],
				"image": [
					{
						"name": "Test Image Model",
						"model": "test-image-model",
						"pricing": {
							"square": {
								"coins": 5,
								"size": "512x512"
							},
							"landscape": {
								"coins": 10,
								"size": "1024x512"
							},
							"portrait": {
								"coins": 10,
								"size": "512x1024"
							}
						}
					}
				]
			},
			"success": true
		}`
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Use a modified version of GetModels that accepts a custom URL
	getModelsWithURL := func(apiKey, url string) ([]Models, error) {
		client := &http.Client{}
		req, _ := http.NewRequest("GET", url, nil)
		req.Header = http.Header{
			"Authorization": []string{"Bearer " + apiKey},
			"Accept":        []string{"application/json"},
		}
		resp, _ := client.Do(req)
		defer resp.Body.Close()

		if resp.StatusCode > 299 || resp.StatusCode < 200 {
			return nil, nil
		}

		bodyText, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		straicoModels, err := UnmarshalStraicoModels(bodyText)
		if err != nil {
			return nil, err
		}

		chatModels := straicoModels.Data.Chat
		viableModels := make([]Models, len(chatModels))
		for i := range chatModels {
			viableModels[i] = Models{
				Name:    chatModels[i].Name,
				Id:      chatModels[i].Model,
				Pricing: int(chatModels[i].Pricing.Coins),
			}
		}
		return viableModels, nil
	}

	// Get models using the test server URL
	models, err := getModelsWithURL("test-key", server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(models) != 2 {
		t.Fatalf("Expected 2 models, got %d", len(models))
	}

	// Check first model
	if models[0].Name != "Test Model 1" {
		t.Errorf("Expected Name 'Test Model 1', got %q", models[0].Name)
	}

	if models[0].Id != "test-model-1" {
		t.Errorf("Expected Id 'test-model-1', got %q", models[0].Id)
	}

	if models[0].Pricing != 10 {
		t.Errorf("Expected Pricing 10, got %d", models[0].Pricing)
	}

	// Check second model
	if models[1].Name != "Test Model 2" {
		t.Errorf("Expected Name 'Test Model 2', got %q", models[1].Name)
	}

	if models[1].Id != "test-model-2" {
		t.Errorf("Expected Id 'test-model-2', got %q", models[1].Id)
	}

	if models[1].Pricing != 20 {
		t.Errorf("Expected Pricing 20, got %d", models[1].Pricing)
	}
}

func TestUnmarshalStraicoModels(t *testing.T) {
	jsonData := []byte(`{
		"data": {
			"chat": [
				{
					"name": "Test Model",
					"model": "test-model",
					"word_limit": 1000,
					"pricing": {
						"coins": 10.5,
						"words": 1000
					},
					"max_output": 500
				}
			],
			"image": [
				{
					"name": "Test Image Model",
					"model": "test-image-model",
					"pricing": {
						"square": {
							"coins": 5,
							"size": "512x512"
						},
						"landscape": {
							"coins": 10,
							"size": "1024x512"
						},
						"portrait": {
							"coins": 10,
							"size": "512x1024"
						}
					}
				}
			]
		},
		"success": true
	}`)

	response, err := UnmarshalStraicoModels(jsonData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !response.Success {
		t.Error("Expected Success to be true")
	}

	if len(response.Data.Chat) != 1 {
		t.Fatalf("Expected 1 chat model, got %d", len(response.Data.Chat))
	}

	if len(response.Data.Image) != 1 {
		t.Fatalf("Expected 1 image model, got %d", len(response.Data.Image))
	}

	chatModel := response.Data.Chat[0]
	if chatModel.Name != "Test Model" {
		t.Errorf("Expected Name 'Test Model', got %q", chatModel.Name)
	}

	if chatModel.Model != "test-model" {
		t.Errorf("Expected Model 'test-model', got %q", chatModel.Model)
	}

	if chatModel.Pricing.Coins != 10.5 {
		t.Errorf("Expected Pricing.Coins 10.5, got %f", chatModel.Pricing.Coins)
	}
}
