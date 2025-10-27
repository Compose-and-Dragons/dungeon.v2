package main

import (
	"context"
	"fmt"
	"log"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

func main() {
	// Docker Model Runner Chat base URL
	engineURL := "http://localhost:12434/engines/v1/"
	chatModelId := "philippecharriere494/queen-pedauque:0.5b-0.0.0"

	client := openai.NewClient(
		option.WithBaseURL(engineURL),
		option.WithAPIKey("I💙DockerModelRunner"),
	)

	ctx := context.Background()

	messages := []openai.ChatCompletionMessageParamUnion{
		//openai.SystemMessage("You are an expert of medieval role playing games."),
		openai.UserMessage("Who are you?"),
		//openai.UserMessage("Tell me something about the Chocolatine."),
	}

	param := openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       chatModelId,
		//Temperature: openai.Opt(0.0),
	}

	stream := client.Chat.Completions.NewStreaming(ctx, param)

	for stream.Next() {
		chunk := stream.Current()
		// Stream each chunk as it arrives
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}

	if err := stream.Err(); err != nil {
		log.Fatalln("😡:", err)
	}
}
