package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	howto "github.com/guitaricet/howto/pkg"
)

const VERSION = "2.0.2"

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: howto <prompt>")
		fmt.Println("To use howto, pass it a prompt to complete. For example: " + howto.GetRandomUsageExample())
		fmt.Println("Options:")

		flag.PrintDefaults()
	}

	do_setup := flag.Bool("setup", false, "Set up howto for the first time")
	do_config := flag.Bool("config", false, "Show the current configuration")
	do_change_prompt := flag.Bool("change-prompt", false, "Change the system message prompt")
	flag.Parse()

	if *do_config {
		howto.PrintEnvInfo()
		return
	}
	if *do_change_prompt {
		howto.ChangeSystemMessage()
		return
	}

	config, err := howto.GetConfig()

	config_does_not_exist := os.IsNotExist(err)
	if *do_setup || config_does_not_exist {
		time.Sleep(1 * time.Second)
		howto.Setup(VERSION)
		return
	}

	if err != nil && !config_does_not_exist {
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

	if len(input) == 0 {
		fmt.Println("Usage: howto <prompt>")
		fmt.Println("To use howto, pass it a prompt to complete. For example: " + howto.GetRandomUsageExample())
		return
	}

	var command string
	command, err = howto.GenerateShellCommand(input, config)

	if err != nil {
		fmt.Println("Error generating command: " + err.Error())
		os.Exit(1)
	}

	if len(command) == 0 {
		fmt.Println("Generated command is empty. Please try to rephrase your prompt.")
	}

	fmt.Println(command)
}
