package main

import (
	"context"
	npcagents "dungeon-master/npc-agents"
	"encoding/json"
	"log"
	"strings"

	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/agents"
	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/helpers"
	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/tools"
	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/ui"

	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/plugins/mcp"
)

var agentsTeam map[string]*agents.NPCAgent
var selectedAgent *agents.NPCAgent

func main() {

	ctx := context.Background()

	llmURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	mcpHost := helpers.GetEnvOrDefault("MCP_SERVER_BASE_URL", "http://localhost:9011/mcp")
	dungeonMasterModel := "openai/" + helpers.GetEnvOrDefault("DUNGEON_MASTER_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")

	fmt.Println("🌍 LLM URL:", llmURL)
	fmt.Println("🌍 MCP Host:", mcpHost)
	fmt.Println("🌍 Dungeon Master Model:", dungeonMasterModel)

	// similaritySearchLimit := helpers.StringToFloat(helpers.GetEnvOrDefault("SIMILARITY_LIMIT", "0.5"))
	// similaritySearchMaxResults := helpers.StringToInt(helpers.GetEnvOrDefault("SIMILARITY_MAX_RESULTS", "2"))

	// [MCP Client] to connect to the [MCP Dungeon Server]
	mcpClient, err := mcp.NewGenkitMCPClient(mcp.MCPClientOptions{
		Name: "c&d",
		StreamableHTTP: &mcp.StreamableHTTPConfig{
			BaseURL: helpers.GetEnvOrDefault("MCP_SERVER_BASE_URL", "http://localhost:9011/mcp"), // docker-mcp-gateway
		},
	})
	if err != nil {
		panic(fmt.Errorf("failed to create MCP client: %v", err))
	}

	ui.Println(ui.Orange, "MCP Client initialized successfully")

	// ---------------------------------------------------------
	// Get the [MCP Tools Index] from the [MCP Client]
	// ---------------------------------------------------------
	toolsRefs := tools.MCPCatalog(ctx, mcpClient)

	// ---------------------------------------------------------
	// AGENT: This is the Dungeon Master Agent using tools
	// ---------------------------------------------------------
	dungeonMasterToolsAgentName := helpers.GetEnvOrDefault("DUNGEON_MASTER_NAME", "Sam")
	dungeonMasterModeltemperature := helpers.StringToFloat(helpers.GetEnvOrDefault("DUNGEON_MASTER_MODEL_TEMPERATURE", "0.0"))
	dungeonMasterModeltopP := helpers.StringToFloat(helpers.GetEnvOrDefault("DUNGEON_MASTER_MODEL_TOP_P", "0.9"))

	dungeonMasterConfig := agents.Config{
		EngineURL:    llmURL,
		Temperature:  dungeonMasterModeltemperature,
		TopP:         dungeonMasterModeltopP,
		ChatModelId:  dungeonMasterModel,
		ToolsModelId: dungeonMasterModel,
		Tools:        toolsRefs,
	}

	// SYSTEM MESSAGE:
	instructions := fmt.Sprintf(`Your name is "%s the Dungeon Master".`, dungeonMasterToolsAgentName) + "\n" + helpers.GetEnvOrDefault("DUNGEON_MASTER_SYSTEM_INSTRUCTIONS", dungeonMasterToolsAgentName)

	dungeonMasterToolsAgent := &agents.NPCAgent{}
	dungeonMasterToolsAgent.Initialize(ctx, dungeonMasterConfig, dungeonMasterToolsAgentName)

	dungeonMasterToolsAgent.SetSystemInstructions(instructions)

	// ---------------------------------------------------------
	// AGENTS:
	// ---------------------------------------------------------
	guardAgent := npcagents.GetGuardAgent(ctx)
	sorcererAgent := npcagents.GetSorcererAgent(ctx)
	healerAgent := npcagents.GetHealerAgent(ctx)
	merchantAgent := npcagents.GetMerchantAgent(ctx)

	// ---------------------------------------------------------
	// [REMOTE] AGENT: This is the Boss agent
	// ---------------------------------------------------------
	// [TODO]

	// ---------------------------------------------------------
	// AGENTS: Creating the Agents Team of the Dungeon
	// ---------------------------------------------------------
	idDungeonMasterToolsAgent := strings.ToLower(dungeonMasterToolsAgentName)
	idGuardAgent := strings.ToLower(guardAgent.Name)
	idSorcererAgent := strings.ToLower(sorcererAgent.Name)
	idHealerAgent := strings.ToLower(healerAgent.Name)
	idMerchantAgent := strings.ToLower(merchantAgent.Name)

	// ---------------------------------------------------------
	// TEAM: Assemble the agents into a team
	// ---------------------------------------------------------
	agentsTeam = map[string]*agents.NPCAgent{
		idDungeonMasterToolsAgent: dungeonMasterToolsAgent,
		idGuardAgent:              guardAgent,
		idSorcererAgent:           sorcererAgent,
		idHealerAgent:             healerAgent,
		idMerchantAgent:           merchantAgent,
	}
	selectedAgent = agentsTeam[idDungeonMasterToolsAgent]

	DisplayAgentsTeam()

	// Loop to interact with the agents
	for {

		var promptText string
		if selectedAgent.Name == dungeonMasterToolsAgentName {
			// Dungeon Master prompt
			promptText = "🤖 (/bye to exit) [" + selectedAgent.Name + "]>"
		} else {
			// Non Player Character prompt
			promptText = "🙂 (/bye to exit /dm to go back to the DM) [" + selectedAgent.Name + "]>"
		}

		// USER PROMPT: (input)
		content, _ := ui.SimplePrompt(promptText, "Type your message here...")

		// ---------------------------------------------------------
		// [COMMAND]: `/bye` command to exit the loop
		// ---------------------------------------------------------
		if strings.HasPrefix(content.Input, "/bye") {
			fmt.Println("👋 Goodbye! Thanks for playing!")
			break
		}

		// ---------------------------------------------------------
		// [COMMAND] `/dm` Back to the Dungeon Master
		// ---------------------------------------------------------
		if strings.HasPrefix(content.Input, "/back-to-dm") || strings.HasPrefix(content.Input, "/dm") || strings.HasPrefix(content.Input, "/dungeonmaster") && selectedAgent.Name != dungeonMasterToolsAgentName {
			selectedAgent = agentsTeam[idDungeonMasterToolsAgent]
			ui.Println(ui.Pink, "👋 You are back to the Dungeon Master:", selectedAgent.Name)
			continue
			/*
				In Go, the continue keyword in a loop immediately jumps to the next iteration of the loop, skipping the rest
				of the code in the current iteration.

				Specifically:
				- In a for loop, continue returns to the beginning of the loop for the next iteration
				- Code after continue in the same iteration is not executed
				- The loop condition is evaluated normally
			*/
		}

		// ---------------------------------------------------------
		// [COMMAND] `/agents` Get the AGENTS team list
		// ---------------------------------------------------------
		if strings.HasPrefix(content.Input, "/agents") {
			DisplayAgentsTeam()
			continue
		}

		// ---------------------------------------------------------
		// [COMMAND] `/tools` Get the TOOLS list
		// ---------------------------------------------------------
		if strings.HasPrefix(content.Input, "/tools") {
			DisplayToolsCatalog(toolsRefs)
			continue
		}
		// ---------------------------------------------------------

		// ---------------------------------------------------------
		// DEBUG:
		if strings.HasPrefix(content.Input, "/memory") {
			selectedAgent.DisplayHistory()
			continue
		}

		switch selectedAgent.Name {
		// ---------------------------------------------------------
		// TALK TO: AGENT: **Dungeon Master** [COMPLETION] with [TOOLS]
		// ---------------------------------------------------------
		case dungeonMasterToolsAgentName: // Zephyr the Dungeon Master

			ui.Println(ui.Yellow, "<", selectedAgent.Name, "speaking...>")

			toolCallsResult, err := selectedAgent.DetectAndExecuteToolCallsWithConfirmation(ctx, dungeonMasterConfig, content.Input)
			if err != nil {
				log.Fatal("😡:", err)
			}
			if toolCallsResult.TotalCalls == 0 {
				// [TODO]
			}
			if toolCallsResult.TotalCalls > 0 {

				toolName, value := GetResultOfToolCall(toolCallsResult)
				ui.Println(ui.Green, "🛠️ Tool called:", toolName)
				ui.Println(ui.Blue, "🛠️ Text:", value)

				switch toolName {

				case "c&d_speak_to_somebody":
					// Switch to the selected agent
					var answer struct {
						Name string `json:"name"`
					}
					err = json.Unmarshal([]byte(value), &answer)

					agentId := strings.ToLower(strings.TrimSpace(answer.Name))
					agent, exists := agentsTeam[agentId]

					// [TODO] check if you are in the same room as the NPC

					if exists {
						selectedAgent = agent
						ui.Printf(ui.Pink, "👋 You are now speaking to %s.\n", selectedAgent.Name)
						continue
					} else {
						ui.Printf(ui.Red, "❌ Agent %q not found. Staying with %s.\n", value, selectedAgent.Name)
					}
				}

			}
			/*
				- c&d_collect_gold
				- c&d_collect_magic_potion
				- c&d_create_player
				- c&d_fight_monster
				- c&d_get_current_room_info
				- c&d_get_dungeon_info
				- c&d_get_dungeon_map
				- c&d_get_player_info
				- c&d_is_player_in_same_room_as_npc
				- c&d_move_by_direction
				- c&d_move_player
				- c&d_speak_to_somebody

			*/

			_, err = selectedAgent.StreamCompletion(ctx, dungeonMasterConfig, toolCallsResult.LastMessage, func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
				fmt.Print(chunk.Text())
				return nil
			})

			if err != nil {
				log.Fatal("😡:", err)
			}
			selectedAgent.ResetMessages() // Clear history after each interaction to avoid tool call accumulation

		// ---------------------------------------------------------
		// TALK TO: AGENT:: **GUARD** + [RAG]
		// ---------------------------------------------------------
		case guardAgent.Name:

			ui.Println(ui.Brown, "<", selectedAgent.Name, "speaking...>")

			_, err = selectedAgent.StreamCompletionWithSimilaritySearch(ctx, npcagents.GetGuardAgentConfig(), content.Input, func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
				fmt.Print(chunk.Text())
				return nil
			})

			if err != nil {
				ui.Println(ui.Red, "Error:", err)
			}

		// ---------------------------------------------------------
		// TALK TO: AGENT:: **SORCERER** + [RAG]
		// ---------------------------------------------------------
		case sorcererAgent.Name:

			ui.Println(ui.Purple, "<", selectedAgent.Name, "speaking...>")

			_, err = selectedAgent.StreamCompletionWithSimilaritySearch(ctx, npcagents.GetSorcererAgentConfig(), content.Input, func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
				fmt.Print(chunk.Text())
				return nil
			})

			if err != nil {
				ui.Println(ui.Red, "Error:", err)
			}

		// ---------------------------------------------------------
		// TALK TO: AGENT:: **HEALER** + [RAG]
		// ---------------------------------------------------------
		case healerAgent.Name:

			ui.Println(ui.Magenta, "<", selectedAgent.Name, "speaking...>")

			_, err = selectedAgent.StreamCompletionWithSimilaritySearch(ctx, npcagents.GetHealerAgentConfig(), content.Input, func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
				fmt.Print(chunk.Text())
				return nil
			})

			if err != nil {
				ui.Println(ui.Red, "Error:", err)
			}

		// ---------------------------------------------------------
		// TALK TO: AGENT:: **MERCHANT** + [RAG]
		// ---------------------------------------------------------
		case merchantAgent.Name:

			ui.Println(ui.Cyan, "<", selectedAgent.Name, "speaking...>")

			_, err = selectedAgent.StreamCompletionWithSimilaritySearch(ctx, npcagents.GetMerchantAgentConfig(), content.Input, func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
				fmt.Print(chunk.Text())
				return nil
			})

			if err != nil {
				ui.Println(ui.Red, "Error:", err)
			}

		// ---------------------------------------------------------
		// TALK TO: AGENT:: **BOSS**
		// ---------------------------------------------------------


		default:
			ui.Printf(ui.Cyan, "\n🤖 %s is thinking...\n", selectedAgent.Name)

		}

	}

}

// ---------------------------------------------------------
// HELPERS:
// ---------------------------------------------------------

func DisplayAgentsTeam() {
	for agentId, agent := range agentsTeam {
		ui.Printf(ui.Cyan, "Agent ID: %s agent name: %s\n", agentId, agent.Name)
	}
	fmt.Println()
}

func DisplayToolsCatalog(tools []ai.ToolRef) {
	ui.Println(ui.Green, "📦 Available Tools:")
	for _, tool := range tools {
		fmt.Println("   -", tool.Name())
	}
	fmt.Println()
}

func GetResultOfToolCall(toolCallsResult *agents.ToolCallsResult) (string, string) {
	toolCalled := toolCallsResult.Results[0]["tool_name"]
	result := toolCallsResult.Results[0]["tool_output"]
	resultText := ""
	// Extract content from the result map
	if resultMap, ok := result.(map[string]any); ok {
		if content, ok := resultMap["content"]; ok {
			// content is an array, get the first element
			if contentArray, ok := content.([]any); ok && len(contentArray) > 0 {
				if contentItem, ok := contentArray[0].(map[string]any); ok {
					if text, ok := contentItem["text"]; ok {
						resultText = text.(string)
					}
				}
			}
		}
	}
	return toolCalled.(string), resultText
}
