package npcagents

import (
	"context"
	"log"

	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/agents"
	"github.com/micro-agent/micro-agent-go/agent/helpers"
)

var sorcererAgentConfig agents.Config

func GetSorcererAgentConfig() agents.Config {
	return sorcererAgentConfig
}

func GetSorcererAgent(ctx context.Context) *agents.NPCAgent {
	engineURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/v1/")
	chatModelId := "openai/" + helpers.GetEnvOrDefault("SORCERER_MODEL", "ai/qwen2.5:1.5B-F16")
	embeddingsModelId := helpers.GetEnvOrDefault("EMBEDDING_MODEL", "ai/mxbai-embed-large:latest")

	agentName := helpers.GetEnvOrDefault("SORCERER_NAME", "Dewey")
	temperature := helpers.StringToFloat(helpers.GetEnvOrDefault("SORCERER_MODEL_TEMPERATURE", "0.0"))
	topP := helpers.StringToFloat(helpers.GetEnvOrDefault("SORCERER_MODEL_TOP_P", "0.9"))

	similaritySearchLimit := helpers.StringToFloat(helpers.GetEnvOrDefault("SIMILARITY_LIMIT", "0.5"))
	similaritySearchMaxResults := helpers.StringToInt(helpers.GetEnvOrDefault("SIMILARITY_MAX_RESULTS", "2"))

	sorcererAgentConfig = agents.Config{
		EngineURL:                  engineURL,
		SimilaritySearchLimit:      similaritySearchLimit,
		SimilaritySearchMaxResults: similaritySearchMaxResults,
		Temperature:                temperature,
		TopP:                       topP,
		ChatModelId:                chatModelId,
		EmbeddingsModelId:          embeddingsModelId,
	}

	sorcererAgent := &agents.NPCAgent{}
	sorcererAgent.Initialize(ctx, sorcererAgentConfig, agentName)

	// ---------------------------------------------------------
	// System Instructions
	// ---------------------------------------------------------
	systemInstructionsPath := helpers.GetEnvOrDefault("SORCERER_SYSTEM_INSTRUCTIONS_PATH", "")
	// [SYSTEM INSTRUCTIONS] Load system instructions from file
	err := sorcererAgent.SetSystemInstructionsFromFile(systemInstructionsPath)
	if err != nil {
		log.Fatal("😡:", err)
	}

	backgroundAndPersonalityPath := helpers.GetEnvOrDefault("SORCERER_CONTEXT_PATH", "")
	// [RAG] Initialize vector store from file(s)
	err = sorcererAgent.InitializeVectorStoreFromFile(ctx, sorcererAgentConfig, backgroundAndPersonalityPath)
	if err != nil {
		log.Fatal("😡:", err)
	}

	return sorcererAgent
}
