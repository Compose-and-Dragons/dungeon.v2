package npcagents

import (
	"context"
	"log"

	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/agents"
	"github.com/micro-agent/micro-agent-go/agent/helpers"
)

var merchantAgentConfig agents.Config

func GetMerchantAgentConfig() agents.Config {
	return merchantAgentConfig
}

func GetMerchantAgent(ctx context.Context) *agents.NPCAgent {
	engineURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/v1/")
	chatModelId := "openai/" + helpers.GetEnvOrDefault("MERCHANT_MODEL", "ai/qwen2.5:1.5B-F16")
	embeddingsModelId := helpers.GetEnvOrDefault("EMBEDDING_MODEL", "ai/mxbai-embed-large:latest")

	agentName := helpers.GetEnvOrDefault("MERCHANT_NAME", "Thorin")
	temperature := helpers.StringToFloat(helpers.GetEnvOrDefault("MERCHANT_MODEL_TEMPERATURE", "0.0"))
	topP := helpers.StringToFloat(helpers.GetEnvOrDefault("MERCHANT_MODEL_TOP_P", "0.9"))

	similaritySearchLimit := helpers.StringToFloat(helpers.GetEnvOrDefault("SIMILARITY_LIMIT", "0.5"))
	similaritySearchMaxResults := helpers.StringToInt(helpers.GetEnvOrDefault("SIMILARITY_MAX_RESULTS", "2"))

	merchantAgentConfig = agents.Config{
		EngineURL:                  engineURL,
		SimilaritySearchLimit:      similaritySearchLimit,
		SimilaritySearchMaxResults: similaritySearchMaxResults,
		Temperature:                temperature,
		TopP:                       topP,
		ChatModelId:                chatModelId,
		EmbeddingsModelId:          embeddingsModelId,
	}

	merchantAgent := &agents.NPCAgent{}
	merchantAgent.Initialize(ctx, merchantAgentConfig, agentName)

	// ---------------------------------------------------------
	// System Instructions
	// ---------------------------------------------------------
	systemInstructionsPath := helpers.GetEnvOrDefault("MERCHANT_SYSTEM_INSTRUCTIONS_PATH", "")
	// [SYSTEM INSTRUCTIONS] Load system instructions from file
	err := merchantAgent.SetSystemInstructionsFromFile(systemInstructionsPath)
	if err != nil {
		log.Fatal("😡:", err)
	}

	backgroundAndPersonalityPath := helpers.GetEnvOrDefault("MERCHANT_CONTEXT_PATH", "")
	// [RAG] Initialize vector store from file(s)
	err = merchantAgent.InitializeVectorStoreFromFile(ctx, merchantAgentConfig, backgroundAndPersonalityPath)
	if err != nil {
		log.Fatal("😡:", err)
	}

	return merchantAgent
}
