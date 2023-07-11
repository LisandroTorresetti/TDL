package news

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

const apiKeyEnvVar = "NEWS_DATA_API_KEY"

var httpClient = &http.Client{}

type Provider interface {
	SummarizeNews(string) (string, error)
}

type New struct {
	Title string `json:"title"`
	Body  string `json:"content"`
	Url   string `json:"link"`
}

type response struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Results      []New  `json:"results"`
}

func GetNew(topic string) (*New, error) {
	// Source: https://newsdata.io/documentation

	newsDataApiKey := os.Getenv(apiKeyEnvVar)

	url := fmt.Sprintf("https://newsdata.io/api/1/news?apikey=%s&language=es", newsDataApiKey)
	if topic != "" {
		url += fmt.Sprintf("&category=%s", topic)
	}

	fmt.Println(url)

	resp, err := httpClient.Get(url)

	if err != nil {
		log.Println("GetNew error -> Cannot make HTTP request: " + err.Error())
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println("GetNew error -> Cannot read body of response: " + err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	var result response
	err = json.Unmarshal(body, &result)

	if err != nil {
		log.Println("GetNew error -> Cannot Unmarshal body of response: " + err.Error())
		return nil, err
	}
	if result.TotalResults < 1 {
		log.Infof("couldn't get news for wanted topic")
		return nil, nil
	}

	return &result.Results[0], nil
}

func GetSummarizedMessage(new *New, gpt Provider) string {
	summarizedBody, err := gpt.SummarizeNews(new.Body)
	if err != nil {
		log.Errorf("couldn't get summarized news: %+v", err)
	}
	return fmt.Sprintf("*%s*\n\n%s\n", new.Title, summarizedBody)
}
