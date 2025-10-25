package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"dungeon-mcp-server/data"
	"dungeon-mcp-server/tools"
	"dungeon-mcp-server/types"

	"github.com/mark3labs/mcp-go/server"

	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/agents"
	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/helpers"
)

func main() {

	// ---------------------------------------------------------
	// NOTE: Create [MCP Server]
	// ---------------------------------------------------------
	s := server.NewMCPServer(
		"dungeon-mcp-server",
		"0.0.0",
	)

	// ---------------------------------------------------------
	// Create a "micro" agent
	// ---------------------------------------------------------
	ctx := context.Background()

	fmt.Println("🤖 Initializing Dungeon Agent...")
	baseURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	fmt.Println("🌍 Model Runner Base URL:", baseURL)
	dungeonModel := helpers.GetEnvOrDefault("DUNGEON_MODEL", "ai/qwen2.5:1.5B-F16")
	fmt.Println("🧠 Dungeon Model:", dungeonModel)

	temperature := helpers.StringToFloat(helpers.GetEnvOrDefault("DUNGEON_MODEL_TEMPERATURE", "0.7"))

	// NOTE: [Agent] Creation
	config := agents.Config{
		EngineURL:   baseURL,
		Temperature: temperature,
		ChatModelId: "openai/" + dungeonModel,
	}

	dungeonAgent := agents.NPCAgent{}
	dungeonAgent.Initialize(ctx, config, "dungeon-agent")

	// ---------------------------------------------------------
	// Game initialisation
	// ---------------------------------------------------------

	// NOTE: Initialize the Player struct
	currentPlayer := types.Player{
		Name: "Unknown",
	}

	width := helpers.StringToInt(helpers.GetEnvOrDefault("DUNGEON_WIDTH", "3"))
	height := helpers.StringToInt(helpers.GetEnvOrDefault("DUNGEON_HEIGHT", "3"))
	entranceX := helpers.StringToInt(helpers.GetEnvOrDefault("DUNGEON_ENTRANCE_X", "0"))
	entranceY := helpers.StringToInt(helpers.GetEnvOrDefault("DUNGEON_ENTRANCE_Y", "0"))
	exitX := helpers.StringToInt(helpers.GetEnvOrDefault("DUNGEON_EXIT_X", "2"))
	exitY := helpers.StringToInt(helpers.GetEnvOrDefault("DUNGEON_EXIT_Y", "2"))

	dungeonName := helpers.GetEnvOrDefault("DUNGEON_NAME", "The Dark Labyrinth")
	dungeonDescription := helpers.GetEnvOrDefault("DUNGEON_DESCRIPTION", "A sprawling underground maze filled with monsters, traps, and treasure.")

	fmt.Println("🧙 Dungeon Name:", dungeonName)
	fmt.Println("📝 Dungeon Description:", dungeonDescription)

	fmt.Println("🏰 Dungeon Size:", width, "x", height)

	// NOTE: Initialize the Dungeon structure
	dungeon := types.Dungeon{
		Name:        dungeonName,
		Description: dungeonDescription,
		Width:       width,
		Height:      height,
		Rooms:       []types.Room{},
		EntranceCoords: types.Coordinates{
			X: entranceX,
			Y: entranceY,
		},
		ExitCoords: types.Coordinates{
			X: exitX,
			Y: exitY,
		},
	}

	// make the dungeon settings configurable via env vars or a config file
	fmt.Println("🚪 Dungeon Entrance Coords:", dungeon.EntranceCoords)
	fmt.Println("🚪 Dungeon Exit Coords:", dungeon.ExitCoords)

	// Create the entrance room of the dungeon

	// ---------------------------------------------------------
	// BEGIN: Generate the entrance room with the dungeon agent
	// ---------------------------------------------------------
	dungeonAgentRoomSystemInstruction := helpers.GetEnvOrDefault("DUNGEON_AGENT_ROOM_SYSTEM_INSTRUCTION", "You are a Dungeon Master. You create rooms in a dungeon. Each room has a name and a short description.")
	dungeonAgent.SetSystemInstructions(dungeonAgentRoomSystemInstruction)

	response, err := dungeonAgent.JsonCompletion(ctx, config, data.Room{}, `
		Create an dungeon entrance room with a name and a short description.
	`)

	if err != nil {
		fmt.Println("🔴 Error generating room:", err)
		return
	}

	fmt.Println("📝 Dungeon Entrance Room Response:", response)

	var roomResponse struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err = json.Unmarshal([]byte(response), &roomResponse); err != nil {
		fmt.Println("Error unmarshaling room response:", err)
		return
	}

	fmt.Println("👋🏰 Entrance Room:", roomResponse)
	// ---------------------------------------------------------
	// END: of Generate the entrance room with the dungeon agent
	// ---------------------------------------------------------
	// NOTE: Initialize the Room structure
	entranceRoom := types.Room{
		ID:          "room_0_0",
		Name:        roomResponse.Name,
		Description: roomResponse.Description,
		IsEntrance:  true,
		IsExit:      false,
		Coordinates: types.Coordinates{
			X: entranceX,
			Y: entranceY,
		},
		Visited:               true,
		HasMonster:            false,
		HasNonPlayerCharacter: false,
		HasTreasure:           false,
		HasMagicPotion:        false,
	}
	dungeon.Rooms = append(dungeon.Rooms, entranceRoom)

	// ---------------------------------------------------------
	// TOOLS Registration
	// ---------------------------------------------------------
	// ---------------------------------------------------------
	// Register tools and their handlers
	// 🤚 These tools will be used by the dungeon-master program
	// ---------------------------------------------------------
	// Create Player
	createPlayerToolInstance := tools.CreatePlayerTool()
	s.AddTool(createPlayerToolInstance, tools.CreatePlayerToolHandler(&currentPlayer, &dungeon))

	// Get Player Info
	getPlayerInfoToolInstance := tools.GetPlayerInformationTool()
	s.AddTool(getPlayerInfoToolInstance, tools.GetPlayerInformationToolHandler(&currentPlayer, &dungeon))

	// Get Dungeon Info
	getDungeonInfoToolInstance := tools.GetDungeonInformationTool()
	s.AddTool(getDungeonInfoToolInstance, tools.GetDungeonInformationToolHandler(&currentPlayer, &dungeon))

	// Move in the dungeon (two variants with same handler)
	moveIntoTheDungeonToolInstance := tools.GetMoveIntoTheDungeonTool()
	s.AddTool(moveIntoTheDungeonToolInstance, tools.MoveByDirectionToolHandler(&currentPlayer, &dungeon, dungeonAgent, config))

	movePlayerToolInstance := tools.GetMovePlayerTool()
	s.AddTool(movePlayerToolInstance, tools.MoveByDirectionToolHandler(&currentPlayer, &dungeon, dungeonAgent, config))

	// Get Current Room Info
	getCurrentRoomInfoToolInstance := tools.GetCurrentRoomInformationTool()
	s.AddTool(getCurrentRoomInfoToolInstance, tools.GetCurrentRoomInformationToolHandler(&currentPlayer, &dungeon))

	// Get Dungeon Map
	getDungeonMapToolInstance := tools.GetDungeonMapTool()
	s.AddTool(getDungeonMapToolInstance, tools.GetDungeonMapToolHandler(&currentPlayer, &dungeon))

	// Collect Gold
	collectGoldToolInstance := tools.CollectGoldTool()
	s.AddTool(collectGoldToolInstance, tools.CollectGoldToolHandler(&currentPlayer, &dungeon))

	// Collect Magic Potion
	collectMagicPotionToolInstance := tools.CollectMagicPotionTool()
	s.AddTool(collectMagicPotionToolInstance, tools.CollectMagicPotionToolHandler(&currentPlayer, &dungeon))

	// Fight Monster
	fightMonsterToolInstance := tools.FightMonsterTool()
	s.AddTool(fightMonsterToolInstance, tools.FightMonsterToolHandler(&currentPlayer, &dungeon))

	// Check if Player is in the same room as an NPC
	isPlayerInSameRoomAsNPCToolInstance := tools.IsPlayerInSameRoomAsNPCTool()
	s.AddTool(isPlayerInSameRoomAsNPCToolInstance, tools.IsPlayerInSameRoomAsNPCToolHandler(&currentPlayer, &dungeon))

	// ---------------------------------------------------------
	// NOTE: Start the [Streamable HTTP MCP server]
	// ---------------------------------------------------------
	httpPort := helpers.GetEnvOrDefault("MCP_HTTP_PORT", "9090")
	fmt.Println("🌍 MCP HTTP Port:", httpPort)

	log.Println("[Dungeon]MCP StreamableHTTP server is running on port", httpPort)

	// Create a custom mux to handle both MCP and health endpoints
	mux := http.NewServeMux()
	// Add healthcheck endpoint (for Docker MCP Gateway with Docker Compose)
	mux.HandleFunc("/health", healthCheckHandler)
	// Add MCP endpoint
	httpServer := server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/mcp"),
	)
	// Register MCP handler with the mux
	mux.Handle("/mcp", httpServer)
	// Start the HTTP server with custom mux
	log.Fatal(http.ListenAndServe(":"+httpPort, mux))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]any{
		"status": "healthy",
	}
	json.NewEncoder(w).Encode(response)
}
