package howto

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"golang.org/x/term"
)

type HowtoConfig struct {
	Version       string `json:"version"`
	Model         string `json:"model"`
	Shell         string `json:"shell"`
	MaxTokens     int    `json:"max_tokens"`
	SystemMessage string `json:"system_message"`
}

type QuestionOptions struct {
	Question        string
	ValidationRegex string
	Secure          bool
}

func GetRandomUsageExample() string {
	example := examples[rand.Intn(len(examples))]
	return example
}

func GetHowtoDir() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "howto")
	} else {
		return filepath.Join(os.Getenv("HOME"), ".howto")
	}
}

func GetConfigPath() string {
	return filepath.Join(GetHowtoDir(), "config.json")
}

func GetConfig() (HowtoConfig, error) {
	var config HowtoConfig

	configPath := GetConfigPath()

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

func PrintConfig() {
	config, err := GetConfig()
	if err != nil {
		fmt.Println("Can't print config: " + err.Error())
		return
	}
	fmt.Println("Config:")

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatalf("Config JSON marshaling failed: %s", err)
		fmt.Printf("%+v\n", config)
		return
	}
	fmt.Println(string(jsonData))
}

func PrintEnvInfo() {
	fmt.Println("OS: " + runtime.GOOS)
	fmt.Println("Config path: " + GetConfigPath())
	PrintConfig()
}

func AskQuestion(opts QuestionOptions) string {
	for {
		fmt.Print(opts.Question)

		var input string
		if opts.Secure {
			bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				fmt.Println("Reading error: " + err.Error())
				fmt.Println("Please try again")
				continue
			}
			input = string(bytePassword)
			fmt.Println()
		} else {
			reader := bufio.NewReader(os.Stdin)
			input, _ = reader.ReadString('\n')
			input = strings.TrimSpace(input)
		}

		if !isValidResponse(input, opts.ValidationRegex) {
			fmt.Println("\nInvalid choice, please try again. The answer should match: " + opts.ValidationRegex)
			continue
		}

		return input
	}
}

func isValidResponse(input string, regex string) bool {
	if regex == "" {
		return true
	}

	matched, err := regexp.MatchString(regex, input)
	if err != nil {
		fmt.Println("Error matching regex: " + err.Error())
		return false
	}
	return matched
}
