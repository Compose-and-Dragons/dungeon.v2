package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"dungeon-mcp-server/types"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/micro-agent/micro-agent-go/agent/helpers"
)

// This code defines the tool for creating a new player in the dungeon game
// It allows creating a character with a name, class and race
// The player is placed at the dungeon entrance with initial stats

func CreatePlayerTool() mcp.Tool {
	return mcp.NewTool("create_player",
		// DESCRIPTION:
		mcp.WithDescription(`Create a new player. Try: "I'm Bob, the Dwarf Warrior."`),

		// PARAMETERS: name, class, race
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("the name of the player"),
		),
		mcp.WithString("class",
			mcp.Required(),
			mcp.Description("the class of the player, e.g., warrior, mage, rogue"),
		),
		mcp.WithString("race",
			mcp.Required(),
			mcp.Description("the race of the player, e.g., human, elf, dwarf"),
		),
	)
}

func CreatePlayerToolHandler(player *types.Player, dungeon *types.Dungeon) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if player.Name != "Unknown" {
			message := "✋ Player already exists: " + player.Name
			fmt.Println(message)
			return mcp.NewToolResultText(message), fmt.Errorf("player already exists: %s", player.Name)
		}

		args := request.GetArguments()
		name := args["name"].(string)
		class := args["class"].(string)
		race := args["race"].(string)

		fmt.Println("👋:", name, class, race)

		*player = types.Player{
			Name:  name,
			Class: class,
			Race:  race,
			Level: 1,
			Position: types.Coordinates{
				X: dungeon.EntranceCoords.X,
				Y: dungeon.EntranceCoords.Y,
			},
			RoomID:     fmt.Sprintf("room_%d_%d", dungeon.EntranceCoords.X, dungeon.EntranceCoords.Y),
			Health:     helpers.StringToInt(helpers.GetEnvOrDefault("PLAYER_INITIAL_HEALTH", "100")),
			Strength:   helpers.StringToInt(helpers.GetEnvOrDefault("PLAYER_INITIAL_STRENGTH", "10")),
			Experience: helpers.StringToInt(helpers.GetEnvOrDefault("PLAYER_INITIAL_EXPERIENCE", "0")),
			GoldCoins:  helpers.StringToInt(helpers.GetEnvOrDefault("PLAYER_INITIAL_GOLD_COINS", "0")),
		}
		playerJSON, err := json.MarshalIndent(*player, "", "  ")
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(string(playerJSON)), nil
	}
}
