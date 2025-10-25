package main

import (
	"context"
	"log"
	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/helpers"
	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/agents"
)

func main() {
	ctx := context.Background()

	//engineURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	engineURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/v1/")

	chatModelId := "openai/" + helpers.GetEnvOrDefault("CHAT_MODEL", "ai/qwen2.5:1.5B-F16")
	embeddingsModelId := helpers.GetEnvOrDefault("EMBEDDING_MODEL", "ai/mxbai-embed-large")

	agentName := helpers.GetEnvOrDefault("SORCERER_NAME", "Elara")

	similaritySearchLimit := helpers.StringToFloat(helpers.GetEnvOrDefault("SIMILARITY_LIMIT", "0.5"))
	similaritySearchMaxResults := helpers.StringToInt(helpers.GetEnvOrDefault("SIMILARITY_MAX_RESULTS", "2"))

	temperature := helpers.StringToFloat(helpers.GetEnvOrDefault("SORCERER_MODEL_TEMPERATURE", "0.0"))
	topP := helpers.StringToFloat(helpers.GetEnvOrDefault("SORCERER_MODEL_TOP_P", "0.9"))

	config := agents.Config{
		EngineURL:                  engineURL,
		SimilaritySearchLimit:      similaritySearchLimit,
		SimilaritySearchMaxResults: similaritySearchMaxResults,
		Temperature:                temperature,
		TopP:                       topP,
		ChatModelId:                chatModelId,
		EmbeddingsModelId:          embeddingsModelId,
	}

	sorcererAgent := agents.NPCAgent{}
	sorcererAgent.Initialize(ctx, config, agentName)

	// ---------------------------------------------------------
	// System Instructions
	// ---------------------------------------------------------
	systemInstructionsPath := helpers.GetEnvOrDefault("SORCERER_SYSTEM_INSTRUCTIONS_PATH", "")
	err := sorcererAgent.SetSystemInstructionsFromFile(systemInstructionsPath)
	if err != nil {
		log.Fatal("😡:", err)
	}

	backgroundAndPersonalityPath := helpers.GetEnvOrDefault("SORCERER_CONTEXT_PATH", "")
	err = sorcererAgent.InitializeVectorStoreFromFile(ctx, config, backgroundAndPersonalityPath)
	if err != nil {
		log.Fatal("😡:", err)
	}

	sorcererAgent.LoopCompletion(ctx, config)

}
