package main

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/compat_oai/openai"
	"github.com/openai/openai-go/option"
)

func main() {
	ctx := context.Background()

	engineURL := "http://localhost:12434/engines/v1/"
	chatModelId := "openai/philippecharriere494/queen-pedauque:1.5b-0.0.0"

	fmt.Println("🌍 LLM URL:", engineURL)
	fmt.Println("🤖 Chat Model:", chatModelId)

	oaiPlugin := &openai.OpenAI{
		APIKey: "I💙DockerModelRunner",
		Opts: []option.RequestOption{
			option.WithBaseURL(engineURL),
		},
	}
	g := genkit.Init(ctx, genkit.WithPlugins(oaiPlugin))

	_, err := genkit.Generate(ctx, g,
		ai.WithModelName(chatModelId),
		ai.WithPrompt("I love the pains au chocolat"),

		ai.WithStreaming(func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
			fmt.Print(chunk.Text())
			return nil
		}),
	)
	if err != nil {
		fmt.Println("😡 Error during generation:", err)
	}

	fmt.Println()
	fmt.Println()

}
