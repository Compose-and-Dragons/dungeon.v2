package tools

import (
	"context"
	"dungeon-mcp-server/types"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

func GetCurrentRoomInformationTool() mcp.Tool {
	return mcp.NewTool("get_current_room_info",
		mcp.WithDescription(`Get information about the current room where the player is located. Try: "Where am I?" or "Look around"`),
	)
}

func GetCurrentRoomInformationToolHandler(player *types.Player, dungeon *types.Dungeon) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

		if result, err := checkPlayerExists(player); err != nil {
			return result, err
		}

		if player.RoomID == "" {
			message := "🚫 Player is not in any room. Please move into the dungeon first."
			fmt.Println(message)
			return mcp.NewToolResultText(message), fmt.Errorf("player not in any room")
		}

		// Find the current room
		var currentRoom *types.Room
		for i := range dungeon.Rooms {
			if dungeon.Rooms[i].ID == player.RoomID {
				currentRoom = &dungeon.Rooms[i]
				break
			}
		}

		if currentRoom == nil {
			message := fmt.Sprintf("❌ Room with ID '%s' not found in dungeon.", player.RoomID)
			fmt.Println(message)
			return mcp.NewToolResultText(message), fmt.Errorf("room not found")
		}

		// Mark room as visited
		//currentRoom.Visited = true

		roomJSON, err := json.MarshalIndent(*currentRoom, "", "  ")
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(string(roomJSON)), nil
	}
}
