package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

type OpenAiResponse struct {
	Id      string   `json:"intValue"`
	Object  string   `json:"stringValue"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Text string `json:"text"`
}

const VERSION = "1.0.3"

func printEnvInfo() {
	fmt.Println("Howto version: " + VERSION)
	fmt.Println("OS: " + runtime.GOOS)

	httpkey := os.Getenv("OPENAI_API_KEY")
	if httpkey == "" {
		fmt.Println("OpenAI API key: not set")
	} else if httpkey[:3] == "sk-" {
		fmt.Println("OpenAI API key: set")
	} else {
		fmt.Println("OpenAI API key: invalid (does not start with sk-)")
	}

	modelName := os.Getenv("HOWTO_OPENAI_MODEL")
	if modelName == "" {
		fmt.Println("OpenAI model: not set")
	} else {
		fmt.Println("OpenAI model: " + modelName)
	}
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "--help" {
		fmt.Println("Usage: howto <prompt>")
		fmt.Println("To use howto, pass it a prompt to complete. For example: `howto tar file without compression`")
		return
	}
	if len(os.Args) < 2 || os.Args[1] == "--env" {
		printEnvInfo()
		return
	}

	httpkey := os.Getenv("OPENAI_API_KEY")
	if httpkey == "" {
		fmt.Println("OPENAI_API_KEY not set")
		fmt.Println("You can get an API key from https://beta.openai.com/docs/quickstart/add-your-api-key")
		fmt.Println("Once you have an API key, set it in your environment with `export OPENAI_API_KEY=<your key>`")
		os.Exit(1)
	}

	// get env variable HOWTO_OPENAI_MODEL if it exists, else use code-davinci-002
	modelName := os.Getenv("HOWTO_OPENAI_MODEL")
	if modelName == "" {
		modelName = "text-davinci-002"
	}

	input := strings.Join(os.Args[1:], " ")
	prompt := fmt.Sprintf("Bash command to %s:```", input)
	suffix := "```"

	body := []byte(fmt.Sprintf(`{
		"model": "%s",
		"prompt": "%s",
		"suffix": "%s",
		"temperature": 0,
		"max_tokens": 256,
		"top_p": 1,
		"frequency_penalty": 0,
		"presence_penalty": 0
	}`, modelName, prompt, suffix))

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request: ", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+httpkey)

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request: ", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	var openaiResponse OpenAiResponse
	err = json.NewDecoder(resp.Body).Decode(&openaiResponse)
	if err != nil {
		fmt.Println("Error decoding response: ", err)
		os.Exit(1)
	}

	choices := openaiResponse.Choices
	if len(choices) == 0 {
		fmt.Println("OpenAI API disn't respont correctly. Did you correctly set you OPENAI_API_KEY?")
		fmt.Println("Response: ", resp)
		os.Exit(1)
	}

	command := openaiResponse.Choices[0].Text
	// if "```" in command, cut out everything after it
	if index := strings.Index(command, suffix); index != -1 {
		command = command[:index]
	}
	command = strings.Trim(command, "\n")

	fmt.Println(command)
}
