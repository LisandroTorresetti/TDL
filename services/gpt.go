package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var client = &http.Client{}

func Summarize(text string) string {
	payload := strings.NewReader(`{
		"model": "gpt-3.5-turbo",
		"messages": [
			{"role": "system", "content": "You are an assistant used in a Telegram bot that can summarize news. Be brief, I dont want a response with much more than 40 words. Your response should be in the same language that the provided news are."},
			{"role": "user", "content": "`+ text +`"}
		],
		"temperature": 0.2
	}`)
	fmt.Println(payload) // TODO: remove log
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", payload)

	if err != nil {
		fmt.Println("Summarize error: cannot create HTTP POST request", err)
		return ""
	}

	req.Header.Add("Content-Type", "application/json")
	apiKey := os.Getenv("OPENAI_API_KEY")
	req.Header.Add("Authorization", "Bearer " + apiKey)

	resp, err := client.Do(req)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Summarize error: cannot read response body", err)
		return ""
	}

	defer resp.Body.Close()
	fmt.Println(string(body)) // TODO: remove log
	var result map[string]any
	json.Unmarshal([]byte(string(body)), &result)
	summarized := result["choices"].([]any)[0].(map[string]any)["message"].(map[string]any)["content"].(string)

	return summarized
}