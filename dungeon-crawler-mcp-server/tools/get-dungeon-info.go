package tools

import (
	"context"
	"dungeon-mcp-server/types"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
)

func GetDungeonInformationTool() mcp.Tool {
	return mcp.NewTool("get_dungeon_info",
		mcp.WithDescription(`Get the current dungeon's information including its layout, rooms, entrance and exit coordinates.`),
	)
}

func GetDungeonInformationToolHandler(player *types.Player, dungeon *types.Dungeon) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

		if callToolResult, err := checkPlayerExists(player); err != nil {
			return callToolResult, err
		}

		// Create a temporary copy structure with the player information
		var mcpResponse struct {
			Player  types.Player  `json:"player"`
			Dungeon types.Dungeon `json:"dungeon"`
		}

		mcpResponse.Player = *player
		mcpResponse.Dungeon = *dungeon

		dungeonAndPlzyerInformationJSON, err := json.MarshalIndent(mcpResponse, "", "  ")
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(string(dungeonAndPlzyerInformationJSON)), nil
	}
}
