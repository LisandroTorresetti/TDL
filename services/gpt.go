package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const model = "gpt-3.5-turbo"
const maxTokens = 30
const temperature = 0

var client = &http.Client{}

func Summarize() string {
	payload := strings.NewReader(`{
		"model": "gpt-3.5-turbo",
		"messages": [
			{"role": "user", "content": "Felicitame por usar la API usando lenguaje taringuero"}
		],
		"temperature": 0.7
	}`)
	fmt.Println(payload)
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
	
	var result map[string]any
	json.Unmarshal([]byte(string(body)), &result)
	summarized := result["choices"].([]any)[0].(map[string]any)["message"].(map[string]any)["content"].(string)

	return summarized
}