package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var httpClient = &http.Client{}

type New struct {
	title string
	body string
	url string
}

func GetNew(topic string) (*New, error) {
	// Source: https://newsdata.io/documentation

	var newsDataApiKey = os.Getenv("NEWS_DATA_API_KEY")

	url := fmt.Sprintf("https://newsdata.io/api/1/news?apikey=%s&category=%s&language=es", newsDataApiKey, topic)

	fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	
	if err != nil {
		log.Println("GetNew error -> Cannot create HTTP request: " + err.Error())
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	
	resp, err := client.Do(req)

	if err != nil {
		log.Println("GetNew error -> Cannot make HTTP request: " + err.Error())
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("GetNew error -> Cannot read body of response: " + err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	var result map[string]any
	err = json.Unmarshal([]byte(string(body)), &result)

	if err != nil {
		log.Println("GetNew error -> Cannot Unmarshal body of response: " + err.Error())
		return nil, err
	}

	firstResult := result["results"].([]any)[0].(map[string]any)

	new := New {
		title: firstResult["title"].(string),
		body: firstResult["content"].(string),
		url: firstResult["link"].(string),
	}

	return &new, nil
}

func GetSummarizedMessage(new *New) string {
	summarizedBody := Summarize(new.body)
	return fmt.Sprintf("*%s*\n\n%s\n\n%s", new.title, summarizedBody, new.url)
}
