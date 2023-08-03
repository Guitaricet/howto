package howto

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/zalando/go-keyring"
)

func SetOpenAiApiKey(apiKey string) error {
	err := keyring.Set(SERVICE_NAME, "openai_api_key", apiKey)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func GetOpenAiApiKey() (string, error) {
	secret, err := keyring.Get(SERVICE_NAME, "openai_api_key")

	if runtime.GOOS == "darwin" {
		// check if it exists at all
		if err == keyring.ErrNotFound {
			fmt.Println("OpenAI API key not found. Please run `howto --setup` to set it in keyring.")
			return "", err
		}
	} else {
		// many issues with keyring on Linux, use env var
		secret = os.Getenv("OPENAI_API_KEY")
		err = nil
	}

	// check if it's valid
	if secret[:3] != "sk-" {
		fmt.Println("OpenAI API key is invalid. Please run `howto config` to set it.")
		return secret, err
	}

	return secret, nil
}
