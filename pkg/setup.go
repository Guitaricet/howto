package howto

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func Setup(version string) {
	fmt.Println("Welcome to howto!")
	configPath := GetConfigPath()
	// check if config file exists
	_, err := os.Stat(configPath)
	config_exists := !os.IsNotExist(err)
	if config_exists {
		fmt.Println("Config file already exists at " + configPath)
		PrintConfig()
		response := AskQuestion(QuestionOptions{
			Question:        "Do you want to overwrite it? (y/n) ",
			ValidationRegex: "y|n",
			Secure:          false,
		})
		if response == "n" {
			fmt.Println("Howto is all set up! Try `howto tar without compression`")
			return
		}
	}

	fmt.Print("Setting up howto")
	if !config_exists {
		fmt.Print(" for the first time")
	}
	fmt.Println("...")

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		fmt.Println("Error creating config directory: " + err.Error())
		fmt.Println("This should never happen. Please report this bug at https://github.com/guitaricet/howto/issues")
		fmt.Println("Please include the following information:")
		fmt.Println("OS: " + runtime.GOOS)
		fmt.Println("Config path: " + configPath)
		fmt.Println("Config dir: " + configDir)
		os.Exit(1)
	}

	openai_api_key := os.Getenv("OPENAI_API_KEY")
	if openai_api_key != "" {
		fmt.Println("Detected OPENAI_API_KEY environment variable.")
		response := AskQuestion(QuestionOptions{
			Question:        "Do you want to use your OPENAI_API_KEY with howto? (y/n) ",
			ValidationRegex: "y|n",
			Secure:          false,
		})
		if response == "y" {
			fmt.Println("Setting OPNEAI_API_KEY in config file")
		} else {
			openai_api_key = AskQuestion(QuestionOptions{
				Question:        "Please enter your OpenAI API key: ",
				ValidationRegex: "sk-[a-zA-Z0-9]{32}",
				Secure:          true,
			})
		}
	}

	shell := AskQuestion(QuestionOptions{
		Question:        "What shell do you use. Just provide the name, e.g. fish (default: bash)? ",
		ValidationRegex: ".+",
		Secure:          false,
	})
	if shell == "" {
		shell = "bash"
	}

	model := AskQuestion(QuestionOptions{
		Question:        "What model do you want to use? (default: gpt-3.5-turbo) ",
		ValidationRegex: "",
		Secure:          false,
	})
	if model == "" {
		model = "gpt-3.5-turbo"
	}

	config := HowtoConfig{
		Model:         model,
		Shell:         shell,
		MaxTokens:     512,
		SystemMessage: "You are a CLI tool that converts user requests to shell commands. E.g. for `bash command to tar file without compression:`, you should reply `tar -cf file.tar file`.",
		OpenAiApiKey:  openai_api_key,
	}

	file, err := os.Create(configPath)
	if err != nil {
		fmt.Println("Error creating config file: " + err.Error())
		os.Exit(1)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(config)
	if err != nil {
		fmt.Println("Error writing config file: " + err.Error())
		os.Exit(1)
	}

	fmt.Printf("\nSetup complete. Try `howto tar without compression`\n\n")
}

func ChangeSystemMessage() {
	new_message := AskQuestion(QuestionOptions{
		Question:        "What do you want the system message to be? ",
		ValidationRegex: ".+",
		Secure:          false,
	})

	config, err := GetConfig()
	if err != nil {
		fmt.Println("Error reading config file: " + err.Error())
		os.Exit(1)
	}
	config.SystemMessage = new_message
	// save
	configPath := GetConfigPath()
	file, err := os.Create(configPath)
	if err != nil {
		fmt.Println("Error creating config file: " + err.Error())
		os.Exit(1)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(config)
	if err != nil {
		fmt.Println("Error writing config file: " + err.Error())
		os.Exit(1)
	}

	fmt.Printf("\nSystem message changed to `%s`\n\n", new_message)
}
