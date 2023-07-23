package howto

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// State is used to keep track of the most recent messages
type HowToState struct {
	Version                string          `json:"version"`
	Conversation           []OpenAiMessage `json:"conversation"`
	LastConversationUpdate time.Time       `json:"lastConversationUpdate"`
}

func GetStatePath() string {
	return filepath.Join(GetHowtoDir(), "state.json")
}

func InitializeState() error {
	statePath := GetStatePath()
	stateDir := filepath.Dir(statePath)
	if err := os.MkdirAll(stateDir, os.ModePerm); err != nil {
		fmt.Println("Error creating state directory: " + err.Error())
		fmt.Println("This should never happen. Please report this bug at https://github.com/guitaricet/howto/issues")
		fmt.Println("Please include the following information:")
		fmt.Println("OS: " + runtime.GOOS)
		fmt.Println("State path: " + statePath)
		fmt.Println("State dir: " + stateDir)
		return fmt.Errorf("error creating state directory: %w", err)
	}
	config, err := GetConfig()
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	state := HowToState{
		Version:                config.Version,
		Conversation:           []OpenAiMessage{},
		LastConversationUpdate: time.Now(),
	}

	err = state.Save(GetStatePath())
	log.Default().Println("Initialized and saved state")
	if err != nil {
		log.Fatalf("Error saving state: %s", err)
		return fmt.Errorf("error saving state: %w", err)
	}
	return err
}

func GetHowtoState() (HowToState, error) {
	var state HowToState
	statePath := GetStatePath()

	file, err := os.Open(statePath)
	if os.IsNotExist(err) {
		InitializeState()
		file, err = os.Open(statePath)
	}
	if err != nil {
		return state, fmt.Errorf("error opening state file at path %s: %w", statePath, err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&state)
	if err != nil {
		return state, fmt.Errorf("error decoding state file at path %s: %w", statePath, err)
	}

	return state, nil
}

func (h HowToState) Save(savePath string) error {
	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("could not create state file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	log.Default().Printf("Saving state: %+v to %s", h, file.Name())
	err = encoder.Encode(h)
	if err != nil {
		return fmt.Errorf("could not encode state: %w", err)
	}
	return nil
}
