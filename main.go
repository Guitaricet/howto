package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// {
//     "id": "cmpl-623HpLk0do5u61ZSmfvqkefqdyc2T",
//     "object": "text_completion",
//     "created": 1665947309,
//     "model": "code-davinci-002",
//     "choices": [
//         {
//             "text": "\nmore code\n\nmore code\n\nmore code\n\nmore code\n\nmore code",
//             "index": 0,
//             "logprobs": null,
//             "finish_reason": "length"
//         }
//     ],
//     "usage": {
//         "prompt_tokens": 4,
//         "completion_tokens": 256,
//         "total_tokens": 260
//     }
// }

type OpenAiResponse struct {
	Id      string   `json:"intValue"`
	Object  string   `json:"stringValue"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Text string `json:"text"`
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "--help" {
		fmt.Println("Usage: howto <prompt>")
		fmt.Println("To use howto, pass it a prompt to complete. For example: `howto tar file without compression`")
		return
	}

	httpkey := os.Getenv("OPENAI_API_KEY")
	if httpkey == "" {
		fmt.Println("OPENAI_API_KEY not set")
		os.Exit(1)
	}

	// get env variable HOWTO_OPENAI_MODEL if it exists, else use code-davinci-002
	modelName := os.Getenv("HOWTO_OPENAI_MODEL")
	if modelName == "" {
		modelName = "text-davinci-002"
	}

	// concatenate all args over spaces
	var input bytes.Buffer
	for i := 1; i < len(os.Args); i++ {
		input.WriteString(os.Args[i])
		input.WriteString(" ")
	}

	prompt := fmt.Sprintf("Bash command to %s:```", input.String())
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
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+httpkey)

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error making request: ", err)
	}

	defer resp.Body.Close()
	var openaiResponse OpenAiResponse
	err = json.NewDecoder(resp.Body).Decode(&openaiResponse)
	if err != nil {
		fmt.Println("Error decoding response: ", err)
	}

	choices := openaiResponse.Choices
	if len(choices) == 0 {
		fmt.Println("OpenAI API disn't respont correctly. Did you correctly set you OPENAI_API_KEY?")
		fmt.Println("Response: ", resp)
		os.Exit(1)
	}

	command := openaiResponse.Choices[0].Text
	// if "```" in command, cut out everything after it
	if index := bytes.Index([]byte(command), []byte("```")); index != -1 {
		command = command[:index]
	}
	command = string(bytes.Trim([]byte(command), "\n"))

	fmt.Println(command)
}
