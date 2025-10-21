package main

import (
	"context"
	"fmt"
	"log"

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

	modelId := "openai/ai/qwen2.5:3B-F16"

	_, err := genkit.Generate(ctx, g,
		ai.WithModelName(modelId),
		
		ai.WithSystem("You are an expert of medieval role playing games."),
		ai.WithPrompt("[Brief] What is a dungeon crawler game?"),

		// ai.WithMessages(
		// 	ai.NewSystemTextMessage("You are the dungeon master of a D&D game."),
		// 	ai.NewUserTextMessage("Generate a D&D NPC Elf name and all its characteristics."),
		// ),
		ai.WithConfig(map[string]any{"temperature": 0.7}),

		ai.WithStreaming(func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
			// Do something with the chunk...
			fmt.Print(chunk.Text())
			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

}
