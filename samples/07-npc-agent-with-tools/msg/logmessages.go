package msg

import (
	"fmt"
	"npc-agent-with-tools/helpers"
)

func Display(text string, something any) {
	shouldIDisplay := helpers.GetEnvOrDefault("LOG_MESSAGES", "true")
	if helpers.StringToBool(shouldIDisplay) {
		fmt.Println(text, something)
	}
}

func displayMessages(messages ...string) {
	for _, msg := range messages {
		fmt.Println(msg)
	}
}

func DisplayEmbeddingsMessages(messages ...string) {
	shouldIDisplay := helpers.GetEnvOrDefault("LOG_EMBEDDINGS_MESSAGES", "true")
	if helpers.StringToBool(shouldIDisplay) {
		displayMessages(messages...)
	}
}

func DisplaySimilarityMessages(messages ...string) {
	shouldIDisplay := helpers.GetEnvOrDefault("LOG_SIMILARITY_MESSAGES", "true")
	if helpers.StringToBool(shouldIDisplay) {
		displayMessages(messages...)
	}
}

func DisplayError(text string, err error) {
	shouldIDisplay := helpers.GetEnvOrDefault("LOG_ERROR_MESSAGES", "true")
	if helpers.StringToBool(shouldIDisplay) {
		fmt.Println(text, err)
	}
}

func DisplayToolMessages(messages ...string) {
	shouldIDisplay := helpers.GetEnvOrDefault("LOG_TOOL_MESSAGES", "true")
	if helpers.StringToBool(shouldIDisplay) {
		displayMessages(messages...)
	}
}
