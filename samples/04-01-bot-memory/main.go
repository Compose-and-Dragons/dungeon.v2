package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/compat_oai/openai"
	"github.com/openai/openai-go/option"
)

func main() {
	ctx := context.Background()
	g := genkit.Init(ctx, genkit.WithPlugins(&openai.OpenAI{
		APIKey: "tada",
		Opts: []option.RequestOption{
			option.WithBaseURL("http://localhost:12434/engines/v1/"),
		},
	}))

	messages := []*ai.Message{}

	agentName := "Elara"

	systemInstructions := fmt.Sprintf(`
		You are an expert of medieval role playing games
		Your name is %s, Weaver of the Arcane
	`, agentName)

	modelId := "openai/ai/qwen2.5:1.5B-F16"

	/* NOTE:
	   - Hello I'm Philippe
	   - Who are you?
	   - My buddy is Guillaume
	   - Who is my buddy?
	*/

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("🤖🧠 [%s](%s) ask me something - /bye to exit> ", agentName, modelId)
		userMessage, _ := reader.ReadString('\n')

		if strings.HasPrefix(userMessage, "/bye") {
			fmt.Println("👋 Bye!")
			break
		}

		if strings.HasPrefix(userMessage, "/history") {
			fmt.Println("📝 Conversation history:")
			for i, msg := range messages {
				// Convert []*ai.Part to string for display
				var parts []string
				for _, part := range msg.Content {
					parts = append(parts, part.Text)
				}
				fmt.Printf("  [%d] %s: %s\n", i, msg.Role, strings.Join(parts, " "))
			}
			continue
		}

		fullResponse, err := genkit.Generate(ctx, g,
			ai.WithModelName(modelId),
			ai.WithSystem(systemInstructions),
			// WithMessages sets the messages. 
			// These messages will be sandwiched between the system and user prompts.
			ai.WithMessages(
				messages...,
			),
			ai.WithPrompt(userMessage),
			ai.WithConfig(map[string]any{
				"temperature": 1.8,
				"top_p":       0.9,
			}),

			ai.WithStreaming(func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
				fmt.Print(chunk.Text())
				return nil
			}),
		)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println() // New line after the response
		// Append user message to history
		messages = append(messages, ai.NewUserTextMessage(strings.TrimSpace(userMessage)))
		// Append assistant response to history
		messages = append(messages, ai.NewModelTextMessage(strings.TrimSpace(fullResponse.Text())))

	}

}
