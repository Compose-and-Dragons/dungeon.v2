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
		APIKey: "I💙DockerModelRunner",
		Opts: []option.RequestOption{
			option.WithBaseURL("http://localhost:12434/engines/v1/"),
		},
	}))

	modelId := "openai/ai/qwen2.5:0.5B-F16"

	_, err := genkit.Generate(ctx, g,
		ai.WithModelName(modelId),

		ai.WithSystem(`
			You are an expert of medieval role playing games
			Your name is Elara, Weaver of the Arcane
		`),
		ai.WithPrompt("Tell me something about you"),

		ai.WithConfig(map[string]any{"temperature": 1.8}), // TODO: check if there is another notation

		ai.WithStreaming(func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
			fmt.Print(chunk.Text())
			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

}
