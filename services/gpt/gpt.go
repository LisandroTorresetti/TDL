package gpt

import (
	"bot-telegram/services/gpt/config"
	"bot-telegram/services/gpt/internal/domain"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	postCompletionURL = "postCompletion"
	openApiKeyEnv     = "OPENAI_API_KEY"
	contentType       = "application/json"
)

type GPT struct {
	config *config.GPTConfig
	client *http.Client
}

func NewGPT(gptConfig *config.GPTConfig, client *http.Client) *GPT {
	return &GPT{
		config: gptConfig,
		client: client,
	}
}

// SummarizeNews returns the userRequest summarized
// ToDo: if we want, we can send a slice of requests
func (gpt *GPT) SummarizeNews(newsToSummarize string) (string, error) {
	userMessage := domain.NewMessage(domain.UserRole, newsToSummarize)
	complainRequest := domain.NewCompletionRequest(gpt.config, []domain.Message{userMessage})
	complainRequestStr, err := domain.GetCompletionRequestAsString(complainRequest)
	if err != nil {
		return "", err // ToDo: add typed error
	}

	endpoint, ok := gpt.config.Endpoints[postCompletionURL]
	if !ok {
		return "", fmt.Errorf("error endpoint does not exists")
	}

	request, err := http.NewRequest(endpoint.Method, endpoint.URL, strings.NewReader(complainRequestStr))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// ToDo: if we need more headers create a package for them. Doing this we can avoid writing them manually
	request.Header.Add("Content-Type", contentType)
	apiKey := os.Getenv(openApiKeyEnv)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	gpt.client.Timeout = 20 * time.Second

	response, err := gpt.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("error doing request: %v", err)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("Summarize error: cannot read response body: %v", err)
	}

	var completionResponse domain.CompletionResponse
	err = json.Unmarshal(body, &completionResponse)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response: %v", err)
	}

	// ToDo: check if it's better to return the struct instead of the message
	return completionResponse.Messages[0].Content, nil
}
