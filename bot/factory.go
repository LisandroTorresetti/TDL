package bot

import (
	"bot-telegram/services/gpt"
	gptConfigLoader "bot-telegram/services/gpt/config"
	"fmt"
	"net/http"
	"time"
)

const timeout = 1

type NewsProvider interface {
	Summarize(string) (string, error)
}

// CreateNewsBot initialize the news bot with all the services that it needs
func CreateNewsBot() (NewsProvider, error) {
	gptConfig, err := gptConfigLoader.LoadConfig()
	if err != nil {
		fmt.Printf("error loading gpt config: %v", err)
		return nil, err
	}

	// ToDo: check if we can do it different
	client := &http.Client{
		Timeout: timeout * time.Second,
	}

	gptService := gpt.NewGPT(gptConfig, client)
	return newNewsBot(gptService), nil
}
