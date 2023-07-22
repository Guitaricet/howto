package howto

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Choice struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	FinishReason string `json:"finish_reason"`
}

type OpenAiResponse struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Choices []Choice `json:"choices"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type OpenAiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func getBodyOpenAI(messages []OpenAiMessage, config HowtoConfig) (string, error) {
	body := map[string]interface{}{
		"model":             config.Model,
		"messages":          messages,
		"temperature":       0,
		"max_tokens":        config.MaxTokens,
		"top_p":             1,
		"frequency_penalty": 0,
		"presence_penalty":  0,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	return string(jsonBody), nil
}

// generateShellCommandAI makes the command via requesting generate from OpenAI
func GenerateShellCommandOpenAI(inputString string, config HowtoConfig) (string, error) {
	prompt := fmt.Sprintf("%s command to %s", config.Shell, inputString)
	messages := []OpenAiMessage{
		{Role: "system", Content: config.SystemMessage},
		{Role: "user", Content: prompt},
	}

	body, err := getBodyOpenAI(messages, config)

	if err != nil {
		fmt.Println("Error creating request body: " + err.Error())
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", strings.NewReader(body))
	if err != nil {
		fmt.Println("Error creating request: ", err)
		return "", err
	}

	api_key, err := GetOpenAiApiKey()
	if err != nil {
		fmt.Println("Error getting OpenAI API key: ", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+api_key)

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request: ", err)
		return "", err
	}

	defer resp.Body.Close()

	var openaiResponse OpenAiResponse
	err = json.NewDecoder(resp.Body).Decode(&openaiResponse)
	if err != nil {
		fmt.Println("Error decoding response: ", err)
		return "", err
	}

	choices := openaiResponse.Choices
	if len(choices) == 0 {
		fmt.Println("OpenAI API didn't respont correctly. Did you correctly set OPENAI_API_KEY?")
		fmt.Println("Response body: ", string(body))
		fmt.Println("Response: ", resp)
		return "", err
	}

	command := openaiResponse.Choices[0].Message.Content
	command = strings.Trim(command, "\n")

	return command, nil
}
