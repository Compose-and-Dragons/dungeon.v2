package tools

import (
	"dungeon-mcp-server/types"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

func checkPlayerExists(player *types.Player) (*mcp.CallToolResult, error) {
	if player == nil || player.Name == "Unknown" {
		message := "✋ No player exists. Please create a player first."
		fmt.Println(message)
		return mcp.NewToolResultText(message), fmt.Errorf("no player exists")
	}
	return nil, nil
}

func checkPlayerIsInARoom(player *types.Player, dungeon *types.Dungeon) (*types.Room, *mcp.CallToolResult, error) {
	var currentRoom *types.Room
	for i := range dungeon.Rooms {
		if dungeon.Rooms[i].ID == player.RoomID {
			currentRoom = &dungeon.Rooms[i]
			break
		}
	}

	if currentRoom == nil {
		message := "❌ Player is not in any room."
		fmt.Println(message)
		return nil, mcp.NewToolResultText(message), fmt.Errorf("player not in any room")
	}
	return currentRoom, nil, nil
}
