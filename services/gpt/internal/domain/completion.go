package domain

import (
	"bot-telegram/services/gpt/config"
	"encoding/json"
)

// CompletionRequest struct that contains the fields needed to do a request against Chat Completions API
type CompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

func NewCompletionRequest(cfg *config.GPTConfig, userMessages []Message) CompletionRequest {
	var messages []Message
	systemMessage := getSystemMessage()
	messages = append(messages, systemMessage)
	messages = append(messages, userMessages...)

	return CompletionRequest{
		Model:       cfg.Model,
		Messages:    messages,
		Temperature: cfg.Temperature,
	}
}

func GetCompletionRequestAsString(completionRequest CompletionRequest) (string, error) {
	jsonBytes, err := json.Marshal(&completionRequest)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// CompletionResponse struct that contains the fields that a response from Chat Completion API has
type CompletionResponse struct {
	Messages []Message `json:"messages"`
}

func (cr *CompletionResponse) UnmarshalJSON(data []byte) error {
	intermediateStructure := struct {
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
	}{}

	err := json.Unmarshal(data, &intermediateStructure)
	if err != nil {
		return err
	}

	for _, choice := range intermediateStructure.Choices {
		cr.Messages = append(cr.Messages, choice.Message)
	}

	return nil
}
