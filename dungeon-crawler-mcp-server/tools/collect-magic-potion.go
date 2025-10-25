package tools

import (
	"context"
	"dungeon-mcp-server/types"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

func CollectMagicPotionTool() mcp.Tool {
	return mcp.NewTool("collect_magic_potion",
		mcp.WithDescription(`Collect magic potions from the current room if available. Try: "Collect the magic potions"`),
	)
}

func CollectMagicPotionToolHandler(player *types.Player, dungeon *types.Dungeon) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if result, err := checkPlayerExists(player); err != nil {
			return result, err
		}

		currentRoom, callToolResult, err := checkPlayerIsInARoom(player, dungeon)
		if err != nil {
			return callToolResult, err
		}

		if !currentRoom.HasMagicPotion || currentRoom.RegenerationHealth <= 0 {
			message := fmt.Sprintf("🧪 There are no magic potions to collect in %s.", currentRoom.Name)
			fmt.Println(message)
			return mcp.NewToolResultText(message), nil
		}

		collectedPotion := currentRoom.RegenerationHealth
		player.Health += collectedPotion
		currentRoom.HasMagicPotion = false
		currentRoom.RegenerationHealth = 0

		message := fmt.Sprintf("🧪 You collected a magic potion from %s! You gained %d health points. Your current health: %d",
			currentRoom.Name, collectedPotion, player.Health)
		fmt.Println(message)
		return mcp.NewToolResultText(message), nil
	}
}
