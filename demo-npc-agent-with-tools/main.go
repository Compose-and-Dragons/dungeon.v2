package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/agents"
	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/helpers"
	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/tools"

	"github.com/firebase/genkit/go/plugins/mcp"
)

func main() {
	ctx := context.Background()

	//engineURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	engineURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/v1/")
	// IMPORTANT: prefix with "openai/" to use the OpenAI plugin TODO: make this automatic
	toolsModelId := "openai/" + helpers.GetEnvOrDefault("TOOLS_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")

	fmt.Println("🌍 LLM URL:", engineURL)
	fmt.Println("🛠️ Tools Model:", toolsModelId)

	agentName := helpers.GetEnvOrDefault("SORCERER_NAME", "Elara")

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
	toolsRefs := tools.MCPCatalog(ctx, mcpClient)

	config := agents.Config{
		EngineURL:    engineURL,
		Temperature:  temperature,
		TopP:         topP,
		ToolsModelId: toolsModelId,
		Tools:        toolsRefs,
	}

	sorcererAgent := agents.NPCAgent{}
	sorcererAgent.Initialize(ctx, config, agentName)

	// Create system message
	systemMsg := `
	You are a helpful D&D assistant that can roll dice and generate character names.
	Use the appropriate tools when asked to roll dice or generate character names.
	`
	sorcererAgent.SetSystemInstructions(systemMsg)

	/*
		Roll 3 dices with 6 faces each.
		Then generate a character name for an elf.
		Finally, roll 2 dices with 8 faces each.
		After that, generate a character name for a dwarf.
	*/

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("🤖🧠 ask me something - /bye to exit> ")
		userMessage, _ := reader.ReadString('\n')

		if strings.HasPrefix(userMessage, "/bye") {
			fmt.Println("👋 Bye!")
			break
		}
		toolCallsResult, err := sorcererAgent.DetectAndExecuteToolCallsWithConfirmation(ctx, config, userMessage)

		if err != nil {
			log.Fatal("😡:", err)
		}
		fmt.Println("🛠️ Total calls:", toolCallsResult.TotalCalls)
		fmt.Println("🛠️ Results:\n", toolCallsResult.Results)
		fmt.Println("🛠️ Final Answer:\n", toolCallsResult.LastMessage)

		fmt.Println(strings.Repeat("=", 50))

	}

	// toolCallsResult, err = sorcererAgent.DetectAndExecuteToolCalls(ctx, config, `
	// 	Say hello to the world.
	// 	Generate a character name for a human.
	// 	Finally, roll 2 dices with 10 faces each.
	// 	Say hello world to Bob Morane.
	// `)
	// if err != nil {
	// 	log.Fatal("😡:", err)
	// }
	// fmt.Println("🛠️ Total calls:", toolCallsResult.TotalCalls)
	// fmt.Println("🛠️ Results:\n", toolCallsResult.Results)
	// fmt.Println("🛠️ Final Answer:\n", toolCallsResult.LastMessage)

}
