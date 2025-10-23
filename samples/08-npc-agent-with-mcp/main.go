package main

import (
	"context"
	"fmt"
	"log"
	"npc-agent-with-mcp/agents"
	"npc-agent-with-mcp/helpers"
	"npc-agent-with-mcp/tools"
	"os"
	"strings"

	"github.com/firebase/genkit/go/plugins/mcp"
)


func main() {
	ctx := context.Background()

	//engineURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	engineURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/v1/")
	// IMPORTANT: prefix with "openai/" to use the OpenAI plugin TODO: make this automatic
	chatModelId := "openai/" + helpers.GetEnvOrDefault("CHAT_MODEL", "ai/qwen2.5:1.5B-F16")
	embeddingsModelId := helpers.GetEnvOrDefault("EMBEDDING_MODEL", "ai/mxbai-embed-large")
	toolsModelId := "openai/" + helpers.GetEnvOrDefault("TOOLS_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")

	fmt.Println("🌍 LLM URL:", engineURL)
	fmt.Println("🤖 Chat Model:", chatModelId)
	fmt.Println("📝 Embeddings Model:", embeddingsModelId)
	fmt.Println("🛠️ Tools Model:", toolsModelId)

	agentName := helpers.GetEnvOrDefault("SORCERER_NAME", "Elara")

	similaritySearchLimit := helpers.StringToFloat(helpers.GetEnvOrDefault("SIMILARITY_LIMIT", "0.5"))
	similaritySearchMaxResults := helpers.StringToInt(helpers.GetEnvOrDefault("SIMILARITY_MAX_RESULTS", "2"))

	temperature := helpers.StringToFloat(helpers.GetEnvOrDefault("SORCERER_MODEL_TEMPERATURE", "0.0"))
	topP := helpers.StringToFloat(helpers.GetEnvOrDefault("SORCERER_MODEL_TOP_P", "0.9"))

	mcpClient, err := mcp.NewGenkitMCPClient(mcp.MCPClientOptions{
		Name: "c&d",
		StreamableHTTP: &mcp.StreamableHTTPConfig{
			BaseURL: helpers.GetEnvOrDefault("MCP_SERVER_BASE_URL", "http://localhost:9011"), // docker-mcp-gateway
		},
	})

	if err != nil {
		fmt.Println("😡 Error connecting Docker MCP Gateway:", err)
		os.Exit(1)
	}

	// Register MCP tools once
	toolsRefs := tools.Catalog(ctx, mcpClient)

	config := agents.Config{
		EngineURL:                  engineURL,
		SimilaritySearchLimit:      similaritySearchLimit,
		SimilaritySearchMaxResults: similaritySearchMaxResults,
		Temperature:                temperature,
		TopP:                       topP,
		ChatModelId:                chatModelId,
		EmbeddingsModelId:          embeddingsModelId,
		ToolsModelId:               toolsModelId,
		Tools:                      toolsRefs,
	}

	sorcererAgent := agents.NPCAgent{}
	sorcererAgent.Initialize(ctx, config, agentName)

	// ---------------------------------------------------------
	// System Instructions
	// ---------------------------------------------------------
	// systemInstructionsPath := helpers.GetEnvOrDefault("SORCERER_SYSTEM_INSTRUCTIONS_PATH", "")
	// err := sorcererAgent.SetSystemInstructionsFromFile(systemInstructionsPath)
	// if err != nil {
	// 	log.Fatal("😡:", err)
	// }

	// Create system message
	systemMsg := `
	You are a helpful D&D assistant that can roll dice and generate character names.
	Use the appropriate tools when asked to roll dice or generate character names.
	`
	sorcererAgent.SetSystemInstructions(systemMsg)

	// backgroundAndPersonalityPath := helpers.GetEnvOrDefault("SORCERER_CONTEXT_PATH", "")
	// err := sorcererAgent.InitializeVectorStoreFromFile(ctx, config, backgroundAndPersonalityPath)
	// if err != nil {
	// 	log.Fatal("😡:", err)
	// }

	toolCallsResult, err := sorcererAgent.DetectAndExecuteToolCallsWithConfirmation(ctx, config, `
		Roll 3 dices with 6 faces each. 
		Then generate a character name for an elf.
		Finally, roll 2 dices with 8 faces each.
		After that, generate a character name for a dwarf.
	`)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println("🛠️ Total calls:", toolCallsResult.TotalCalls)
	fmt.Println("🛠️ Results:\n", toolCallsResult.Results)
	fmt.Println("🛠️ Final Answer:\n", toolCallsResult.LastMessage)

	fmt.Println(strings.Repeat("=", 50))

	toolCallsResult, err = sorcererAgent.DetectAndExecuteToolCalls(ctx, config, `
		Say hello to the world.
		Generate a character name for a human.
		Finally, roll 2 dices with 10 faces each.
		Say hello world to Bob Morane.
	`)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println("🛠️ Total calls:", toolCallsResult.TotalCalls)
	fmt.Println("🛠️ Results:\n", toolCallsResult.Results)
	fmt.Println("🛠️ Final Answer:\n", toolCallsResult.LastMessage)
	//sorcererAgent.LoopCompletion(ctx, config)

}
