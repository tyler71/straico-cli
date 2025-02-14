package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const modelsApi = "https://api.straico.com/v1/models"

type Models struct {
	Name    string
	Id      string
	Pricing int
}

func GetModels(apiKey string) ([]Models, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", modelsApi, nil)
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + apiKey},
		"Accept":        []string{"application/json"},
	}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		errorMessage := fmt.Errorf("request failed. Error: %s", resp.Status)
		return nil, errorMessage
	}
	bodyText, err := io.ReadAll(resp.Body)

	straicoModels, err := UnmarshalStraicoModels(bodyText)
	if err != nil {
		errorMessage := fmt.Errorf("request failed. Error: %w", err)
		return nil, errorMessage
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

func UnmarshalStraicoModels(data []byte) (ModelsResponse, error) {
	var r ModelsResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ModelsResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ModelsResponse struct {
	Data    Data `json:"data"`
	Success bool `json:"success"`
}

type Data struct {
	Chat  []Chat  `json:"chat"`
	Image []Image `json:"image"`
}

type Chat struct {
	Name      string      `json:"name"`
	Model     string      `json:"model"`
	WordLimit int64       `json:"word_limit"`
	Pricing   ChatPricing `json:"pricing"`
	MaxOutput int64       `json:"max_output"`
}

type ChatPricing struct {
	Coins float64 `json:"coins"`
	Words int64   `json:"words"`
}

type Image struct {
	Name    string       `json:"name"`
	Model   string       `json:"model"`
	Pricing ImagePricing `json:"pricing"`
}

type ImagePricing struct {
	Square    Landscape `json:"square"`
	Landscape Landscape `json:"landscape"`
	Portrait  Landscape `json:"portrait"`
}

type Landscape struct {
	Coins int64  `json:"coins"`
	Size  string `json:"size"`
}
