package npcagents

import (
	"context"
	"log"

	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/agents"
	"github.com/micro-agent/micro-agent-go/agent/helpers"
)
var guardAgentConfig agents.Config

func GetGuardAgentConfig() agents.Config {
	return guardAgentConfig
}


func GetGuardAgent(ctx context.Context) *agents.NPCAgent {
	engineURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/v1/")
	chatModelId := "openai/" + helpers.GetEnvOrDefault("GUARD_MODEL", "ai/qwen2.5:1.5B-F16")
	embeddingsModelId := helpers.GetEnvOrDefault("EMBEDDING_MODEL", "ai/mxbai-embed-large:latest")

	agentName := helpers.GetEnvOrDefault("GUARD_NAME", "Huey")
	temperature := helpers.StringToFloat(helpers.GetEnvOrDefault("GUARD_MODEL_TEMPERATURE", "0.0"))
	topP := helpers.StringToFloat(helpers.GetEnvOrDefault("GUARD_MODEL_TOP_P", "0.9"))

	similaritySearchLimit := helpers.StringToFloat(helpers.GetEnvOrDefault("SIMILARITY_LIMIT", "0.5"))
	similaritySearchMaxResults := helpers.StringToInt(helpers.GetEnvOrDefault("SIMILARITY_MAX_RESULTS", "2"))

	guardAgentConfig = agents.Config{
		EngineURL:                  engineURL,
		SimilaritySearchLimit:      similaritySearchLimit,
		SimilaritySearchMaxResults: similaritySearchMaxResults,
		Temperature:                temperature,
		TopP:                       topP,
		ChatModelId:                chatModelId,
		EmbeddingsModelId:          embeddingsModelId,
	}

	guardAgent := &agents.NPCAgent{}
	guardAgent.Initialize(ctx, guardAgentConfig, agentName)

	// ---------------------------------------------------------
	// System Instructions
	// ---------------------------------------------------------
	systemInstructionsPath := helpers.GetEnvOrDefault("GUARD_SYSTEM_INSTRUCTIONS_PATH", "")
	// [SYSTEM INSTRUCTIONS] Load system instructions from file
	err := guardAgent.SetSystemInstructionsFromFile(systemInstructionsPath)
	if err != nil {
		log.Fatal("😡:", err)
	}

	backgroundAndPersonalityPath := helpers.GetEnvOrDefault("GUARD_CONTEXT_PATH", "")
	// [RAG] Initialize vector store from file(s)
	err = guardAgent.InitializeVectorStoreFromFile(ctx, guardAgentConfig, backgroundAndPersonalityPath)
	if err != nil {
		log.Fatal("😡:", err)
	}

	return guardAgent
}
