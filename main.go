package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
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

type HowtoConfig struct {
	Model     string `json:"model"`
	Shell     string `json:"shell"`
	MaxTokens int    `json:"max_tokens"`
}

type CachedResponse struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type HowToState struct {
	Version     string           `json:"version"`
	Cache       []CachedResponse `json:"cache"`
	LastWarning time.Time        `json:"lastWarning"`
}

const VERSION = "1.1.0-dev"
const DEFAULT_CONFIG = `{
	"model": "text-davinci-002",
	"shell": "bash",
	"max_tokens": 256
}`

func getConfigPath() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("APPDATA") + "\\howto\\config.json"
	} else {
		return os.Getenv("HOME") + "/.howto/config.json"
	}
}

func getConfig() (HowtoConfig, error) {
	var config HowtoConfig

	configPath := getConfigPath()
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		return config, err
	}

	file, err := os.Open(configPath)
	if err != nil {
		return config, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func checkEnv() {
	hasError := false
	if os.Getenv("OPENAI_API_KEY") == "" {
		fmt.Println("Please set the OPENAI_API_KEY environment variable")
		fmt.Println("You can get an API key from https://beta.openai.com/docs/quickstart/add-your-api-key")
		fmt.Println("Once you have an API key, set it in your environment with `export OPENAI_API_KEY=<your key>`")
		hasError = true
	}

	if os.Getenv("HOWTO_OPENAI_MODEL") != "" {
		// this is a very anoying message, so let's only print it with 20% probability
		if rand.Intn(100) < 20 {
			fmt.Printf("The HOWTO_OPENAI_MODEL environment variable is deprecated. ")
			fmt.Printf("use the config file %s instead.\n", getConfigPath())
		}
	}

	if hasError {
		os.Exit(1)
	}
}

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

	config, err := getConfig()
	if os.IsNotExist(err) {
		fmt.Println("Config file: not found")
	} else if err != nil {
		fmt.Println("Config path: " + getConfigPath())
		fmt.Println("Error reading config file: " + err.Error())
	} else {
		fmt.Println("Config path: " + getConfigPath())
		fmt.Printf("Config: %v", config)
	}
}

func setup() {
	fmt.Println("First time setup")
	if os.Getenv("OPENAI_API_KEY") == "" {
		fmt.Println("Please set the OPENAI_API_KEY environment variable")
		fmt.Println("You can get an API key from https://beta.openai.com/docs/quickstart/add-your-api-key")
		fmt.Println("Once you have an API key, set it in your environment with `export OPENAI_API_KEY=<your key>`")
		os.Exit(1)
	}

	configPath := getConfigPath()
	fmt.Println("Creating default config at" + configPath)

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		fmt.Println("Error creating config directory: " + err.Error())
		os.Exit(1)
	}

	file, err := os.Create(configPath)
	if err != nil {
		fmt.Println("Error creating config file: " + err.Error())
		os.Exit(1)
	}
	defer file.Close()

	fmt.Println("What shell do you use?")
	fmt.Println("If you don't know, just press enter to use the default (bash)")
	fmt.Println("You can change this later in the config file " + configPath)
	fmt.Println("Options: bash, zsh, fish, powershell")

	var shell string
	fmt.Scanln(&shell)
	if shell == "" {
		shell = "bash"
	}
	shell = strings.ToLower(shell)
	fmt.Println("Setting shell to " + shell)

	var config HowtoConfig
	err = json.Unmarshal([]byte(DEFAULT_CONFIG), &config)
	if err != nil {
		fmt.Println("Error parsing default config: " + err.Error())
		fmt.Println("This is a serious bug, please report it at https://github.com/guitaricet/howto/issues with the following information:")
		fmt.Println("Default config string:\n" + DEFAULT_CONFIG)
		os.Exit(1)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(config)
	if err != nil {
		fmt.Println("Error writing config file: " + err.Error())
		os.Exit(1)
	}

	fmt.Printf("Setup complete. Now you can use howto!\n\n")
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

	_, err := os.Stat(getConfigPath())
	if os.IsNotExist(err) {
		setup()
	}

	checkEnv() // guarantees that OPENAI_API_KEY is set

	httpkey := os.Getenv("OPENAI_API_KEY")

	config, err := getConfig()
	if err == io.EOF {
		setup()
	} else if err != nil {
		fmt.Println("Error reading config file: " + err.Error())
		os.Exit(1)
	}

	input := strings.Join(os.Args[1:], " ")

	// prompt example: "bash command to tar file without compression: ```[insert]```"
	prompt := fmt.Sprintf("%s command to %s:```", config.Shell, input)
	suffix := "```"

	body := []byte(fmt.Sprintf(`{
		"model": "%s",
		"prompt": "%s",
		"suffix": "%s",
		"temperature": 0,
		"max_tokens": %d,
		"top_p": 1,
		"frequency_penalty": 0,
		"presence_penalty": 0
	}`, config.Model, prompt, suffix, config.MaxTokens))

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
		fmt.Println("OpenAI API disn't respont correctly. Did you correctly set OPENAI_API_KEY?")
		// more info about the response
		fmt.Println("Response body: ", string(body))
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
