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

var examples = []string{
	"howto tar without compression",
	"howto oneline install conda",
	"howto du -hs hidden files",
	"howto donwload from gcp bucket",
	"howto pull from upstream",
	"howto push if the only update is the tag",
	"howto get ubuntu version",
	"howto undo make",
	"howto connect to mongo running inside docker",
	"howto check if something is running on my port 27017",
	"howto get user id for user vlialin",
	"howto create user vlialin with user IDs 5030 and GID 4030 and assign them a home directory in /mnt/shared_home",
	"howto tree withiout node_modules",
	"howto 'grep my zsh history and print all examples containing howto (with trailing space)'",
}

type HowtoConfig struct {
	Version       string `json:"version"`
	Model         string `json:"model"`
	Shell         string `json:"shell"`
	MaxTokens     int    `json:"max_tokens"`
	SystemMessage string `json:"system_message"`
	OpenAiApiKey  string `json:"openai_api_key"`
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
	if os.IsNotExist(err) {
		return config, fmt.Errorf("config file %s does not exist", configPath)
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
	if config.OpenAiApiKey == "" {
		config.OpenAiApiKey = "not set"
	} else if config.OpenAiApiKey[:3] != "sk-" {
		config.OpenAiApiKey = "invalid (does not start with sk-)"
	} else {
		config.OpenAiApiKey = config.OpenAiApiKey[:3] + "..."
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
		} else {
			reader := bufio.NewReader(os.Stdin)
			input, _ = reader.ReadString('\n')
			input = strings.TrimSpace(input)
		}

		if !isValidResponse(input, opts.ValidationRegex) {
			fmt.Println("Please answer with one of the following expression: opts.ValidationRegex")
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
