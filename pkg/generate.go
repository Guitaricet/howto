// contains rule-based logic for generating answer
// including caching and conversation tracking
package howto

import (
	"fmt"
	"time"
)

func GenerateShellCommand(command string, config HowtoConfig) (string, error) {
	state, err := GetHowtoState()
	if err != nil {
		return "", fmt.Errorf("error getting state: %w", err)
	}

	messages := []OpenAiMessage{
		{Role: "system", Content: config.SystemMessage},
	}

	time_delta := time.Since(state.LastConversationUpdate)
	if time_delta.Minutes() <= 1 {
		for _, message := range state.Conversation {
			fmt.Printf("%s: %s\n", message.Role, message.Content)
		}

		messages = append(messages, state.Conversation...)
		// we append command, because prompt does not make sense in a conversation-style request
		messages = append(messages, OpenAiMessage{Role: "user", Content: command})
	} else {
		prompt := fmt.Sprintf("%s command to %s", config.Shell, command)
		messages = append(messages, OpenAiMessage{Role: "user", Content: prompt})
	}

	response, err := GenerateShellCommandOpenAI(messages, config)

	messages = append(messages, OpenAiMessage{Role: "assistant", Content: response})

	state.Conversation = messages[1:]
	state.LastConversationUpdate = time.Now()
	state.Save(GetStatePath())

	return response, err
}
