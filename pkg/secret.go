package howto

import (
	"fmt"
	"log"

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

	// check if it exists at all
	if err == keyring.ErrNotFound {
		fmt.Println("OpenAI API key not found. Please run `howto --setup` to set it.")
		return "", err
	}

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	// check if it's valid
	if secret[:3] != "sk-" {
		fmt.Println("OpenAI API key is invalid. Please run `howto config` to set it.")
		return secret, err
	}

	return secret, nil
}
