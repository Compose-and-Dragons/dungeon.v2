package tools

import (
	"context"
	"dungeon-mcp-server/types"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

type IsPlayerInSameRoomAsNPCResponse struct {
	InSameRoom bool   `json:"in_same_room"`
	PlayerRoom string `json:"player_room_id"`
	NPCRoom    string `json:"npc_room_id,omitempty"`
	Message    string `json:"message"`
}

func IsPlayerInSameRoomAsNPCTool() mcp.Tool {
	return mcp.NewTool("is_player_in_same_room_as_npc",
		mcp.WithDescription("Check if the player is in the same room as a specific NPC."),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name of the NPC to check."),
		),
	)
}

func IsPlayerInSameRoomAsNPCToolHandler(player *types.Player, dungeon *types.Dungeon) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		
		if player.Name == "Unknown" {
			response := IsPlayerInSameRoomAsNPCResponse{
				InSameRoom: false,
				PlayerRoom: "",
				Message:    "✋ No player exists. Please create a player first.",
			}
			responseJSON, _ := json.MarshalIndent(response, "", "  ")
			return mcp.NewToolResultText(string(responseJSON)), fmt.Errorf("no player exists")
		}

		args := request.GetArguments()
		npcName := args["name"].(string)

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

			response := IsPlayerInSameRoomAsNPCResponse{
				InSameRoom: false,
				PlayerRoom: "",
				Message:    "✋ Room not found.",
			}
			responseJSON, _ := json.MarshalIndent(response, "", "  ")
			return mcp.NewToolResultText(string(responseJSON)), fmt.Errorf("room not found")
		}

		// Check if there's an NPC in the current room
		if currentRoom.NonPlayerCharacter == nil {
			response := IsPlayerInSameRoomAsNPCResponse{
				InSameRoom: false,
				PlayerRoom: player.RoomID,
				Message:    fmt.Sprintf("❌ No NPC found in room '%s'", currentRoom.ID),
			}
			responseJSON, _ := json.MarshalIndent(response, "", "  ")
			return mcp.NewToolResultText(string(responseJSON)), nil
		}

		if strings.EqualFold(npcName, currentRoom.NonPlayerCharacter.Name) {
			response := IsPlayerInSameRoomAsNPCResponse{
				InSameRoom: true,
				PlayerRoom: player.RoomID,
				NPCRoom:    currentRoom.ID,
				Message:    fmt.Sprintf("✅ Player is in the same room as NPC '%s' (Room: %s)", currentRoom.NonPlayerCharacter.Name, currentRoom.ID),
			}
			responseJSON, _ := json.MarshalIndent(response, "", "  ")
			return mcp.NewToolResultText(string(responseJSON)), nil
		} else {
			response := IsPlayerInSameRoomAsNPCResponse{
				InSameRoom: false,
				PlayerRoom: player.RoomID,
				NPCRoom:    currentRoom.ID,
				Message:    fmt.Sprintf("❌ Player is not in the same room as NPC '%s'. Player is in room '%s' with NPC '%s'", npcName, player.RoomID, currentRoom.NonPlayerCharacter.Name),
			}
			responseJSON, _ := json.MarshalIndent(response, "", "  ")
			return mcp.NewToolResultText(string(responseJSON)), nil
		}
	}
}