package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/helpers"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/compat_oai/openai"
	"github.com/openai/openai-go/option"
)

func main() {
	ctx := context.Background()

	engineURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/v1/")
	chatModelId := "openai/" + helpers.GetEnvOrDefault("CHAT_MODEL", "philippecharriere494/queen-pedauque:0.5b-0.0.0")

	fmt.Println("🌍 LLM URL:", engineURL)
	fmt.Println("🤖 Chat Model:", chatModelId)


	oaiPlugin := &openai.OpenAI{
		APIKey: "I💙DockerModelRunner",
		Opts: []option.RequestOption{
			option.WithBaseURL(engineURL),
		},
	}
	g := genkit.Init(ctx, genkit.WithPlugins(oaiPlugin))

	messages := []*ai.Message{}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("🤖🧠 (%s) ask me something - /bye to exit> ", chatModelId)
		userMessage, _ := reader.ReadString('\n')

		if strings.HasPrefix(userMessage, "/bye") {
			fmt.Println("👋 Bye!")
			break
		}

		fullResponse, err := genkit.Generate(ctx, g,
			ai.WithModelName(chatModelId),
			// WithMessages sets the messages.
			// These messages will be sandwiched between the system and user prompts.
			ai.WithMessages(
				messages...,
			),
			ai.WithPrompt(userMessage),
			// ai.WithConfig(map[string]any{
			// 	"temperature": 0.0,
			// 	"top_p":       0.9,
			// }),

			ai.WithStreaming(func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
				// Do something with the chunk...
				fmt.Print(chunk.Text())
				return nil
			}),
		)
		if err != nil {
			fmt.Println("😡 Error during generation:", err)
			continue
		}

		// Conversation memory
		// Append user message to history
		messages = append(messages, ai.NewUserTextMessage(strings.TrimSpace(userMessage)))
		// Append assistant response to history
		messages = append(messages, ai.NewModelTextMessage(strings.TrimSpace(fullResponse.Text())))

		fmt.Println()
		//fmt.Println()

	}

}
