package prompt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Prompt struct {
	Message     string   `json:"message"`
	Model       []string `json:"models"`
	FileUrls    []string `json:"file_urls,omitempty"`
	YoutubeUrls []string `json:"youtube_urls,omitempty"`
	MaxToken    int      `json:"max_tokens,omitempty"`
	UrlPrefix   string
}

const MaxContextLength = 25

// Request main entrypoint, This requests from the api and returns the response.
func (p Prompt) Request(key string, text string, context []string) (response StraicoResponse, err error) {
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
	req, _ := http.NewRequest("POST", p.UrlPrefix, bytes.NewBuffer(jsonAbc))

	req.Header = http.Header{
		"Authorization": []string{"Bearer " + key},
		"Content-Type":  []string{"application/json"},
		"Accept":        []string{"application/json"},
	}
	resp, err := client.Do(req)
	if err != nil {
		return StraicoResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		errorMessage := fmt.Errorf("request failed. Error: %s", resp.Status)
		return StraicoResponse{}, errorMessage
	}
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		errorMessage := fmt.Errorf("unable to read body. Error: %s", err.Error())
		return StraicoResponse{}, errorMessage
	}

	llmText, err := UnmarshalStraicoResponse(bodyText)
	if err != nil {
		errorMessage := fmt.Errorf("unable to unmarshal body. Error: %s", err.Error())
		return StraicoResponse{}, errorMessage
	}

	return llmText, nil
}

// This will be fed into straico's api
func (p Prompt) Read(pData []byte) (n int, err error) {
	marshalledData, err := json.Marshal(p)
	copy(pData, marshalledData)
	if err != nil {
		errorMessage := fmt.Errorf("failed to marshal prompt. Error: %s", err.Error())
		return 0, errorMessage
	}

	return len(marshalledData), err
}
