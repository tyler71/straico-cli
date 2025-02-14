package prompt

import "encoding/json"

func UnmarshalStraicoResponse(data []byte) (StraicoResponse, error) {
	var r StraicoResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

type StraicoResponse struct {
	Data    Data `json:"data"`
	Success bool `json:"success"`
}

type Data struct {
	OverallPrice OverallPrice        `json:"overall_price"`
	OverallWords OverallPrice        `json:"overall_words"`
	Completions  map[string]LLMModel `json:"completions"`
}

type LLMModel struct {
	Completion LLMCompletion `json:"completion"`
	Price      OverallPrice  `json:"price"`
	Words      OverallPrice  `json:"words"`
}

type LLMCompletion struct {
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int64       `json:"index"`
	Message      Message     `json:"message"`
	FinishReason string      `json:"finish_reason"`
	Logprobs     interface{} `json:"logprobs"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

type OverallPrice struct {
	Input  float64 `json:"input"`
	Output float64 `json:"output"`
	Total  float64 `json:"total"`
}
