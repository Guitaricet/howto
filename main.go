package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	howto "github.com/guitaricet/howto/pkg"
)

const VERSION = "2.0.0-dev"

func main() {
	flag.Usage = func() {
		fmt.Println("Howto Version " + VERSION)
		fmt.Println("Howto is a command line tool that uses OpenAI API to help you write shell commands.")
		fmt.Println("  Usage: howto <prompt>")
		fmt.Println("  For example: " + howto.GetRandomUsageExample())
		fmt.Println("Other commands:")

		flag.PrintDefaults()
	}

	do_setup := flag.Bool("setup", false, "Set up howto for the first time")
	do_config := flag.Bool("config", false, "Show the current configuration")
	do_change_prompt := flag.Bool("change-prompt", false, "Change the system message prompt")
	flag.Parse()

	if *do_setup {
		howto.Setup(VERSION)
		return
	}
	if *do_config {
		howto.PrintEnvInfo()
		return
	}
	if *do_change_prompt {
		howto.ChangeSystemMessage()
		return
	}

	_, err := os.Stat(howto.GetConfigPath())
	if os.IsNotExist(err) {
		fmt.Println("First time setup")
		time.Sleep(1 * time.Second)
		howto.Setup(VERSION)
	}

	config, err := howto.GetConfig()
	if err != nil {
		fmt.Println("Error reading config file: " + err.Error())
		response := howto.AskQuestion(howto.QuestionOptions{
			Question:        "Do you want to delete your config file and run `howto --setup` again? (y/n) ",
			ValidationRegex: "y|n",
			Secure:          false,
		})
		if response == "y" {
			os.Remove(howto.GetConfigPath())
			howto.Setup(VERSION)
		}
	}

	input := strings.Join(os.Args[1:], " ")

	var command string
	command, err = howto.GenerateShellCommandOpenAI(input, config)

	if err != nil {
		fmt.Println("Error generating command: " + err.Error())
		os.Exit(1)
	}

	if len(command) == 0 {
		fmt.Println("Generated command is empty. Please try to rephrase your prompt.")
	}

	fmt.Println(command)
}
