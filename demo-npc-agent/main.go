package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/agents"
	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/helpers"
)

func main() {
	ctx := context.Background()

	//engineURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	engineURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/v1/")
	// IMPORTANT: prefix with "openai/" to use the OpenAI plugin TODO: make this automatic
	chatModelId := "openai/" + helpers.GetEnvOrDefault("CHAT_MODEL", "ai/qwen2.5:1.5B-F16")
	embeddingsModelId := helpers.GetEnvOrDefault("EMBEDDING_MODEL", "ai/mxbai-embed-large")

	fmt.Println("🌍 LLM URL:", engineURL)
	fmt.Println("🤖 Chat Model:", chatModelId)
	fmt.Println("📝 Embeddings Model:", embeddingsModelId)

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
	// [RAG] Initialize vector store from file(s)
	err = sorcererAgent.InitializeVectorStoreFromFile(ctx, config, backgroundAndPersonalityPath)
	if err != nil {
		log.Fatal("😡:", err)
	}

	sorcererAgent.LoopCompletion(ctx, config)

}
