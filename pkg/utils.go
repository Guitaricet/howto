package howto

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"golang.org/x/term"
)

type HowtoConfig struct {
	Version       string `json:"version"`
	Model         string `json:"model"`
	Shell         string `json:"shell"`
	MaxTokens     int    `json:"max_tokens"`
	SystemMessage string `json:"system_message"`
}

type HowToState struct {
	Version      string            `json:"version"`
	Cache        map[string]string `json:"cache"`
	Conversation []string          `json:"conversation"`
	LastWarning  time.Time         `json:"lastWarning"`
}

type QuestionOptions struct {
	Question        string
	ValidationRegex string
	Secure          bool
}

func GetRandomUsageExample() string {
	rand.Seed(time.Now().Unix())
	example := examples[rand.Intn(len(examples))]
	return example
}

func GetConfigPath() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("APPDATA") + "\\howto\\config.json"
	} else {
		return os.Getenv("HOME") + "/.howto/config.json"
	}
}

func GetConfig() (HowtoConfig, error) {
	var config HowtoConfig

	configPath := GetConfigPath()
	_, err := os.Stat(configPath)
	if err != nil {
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
