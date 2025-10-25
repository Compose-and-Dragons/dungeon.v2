package tools

import (
	"context"
	"dungeon-mcp-server/types"
	"fmt"
	"strings"
	"unicode"

	"github.com/mark3labs/mcp-go/mcp"
)

func GetDungeonMapTool() mcp.Tool {
	return mcp.NewTool("get_dungeon_map",
		// DESCRIPTION:
		mcp.WithDescription(`Generate an ASCII map of the discovered dungeon rooms showing the player position, NPCs, and monsters with a legend.`),
	)
}

func GetDungeonMapToolHandler(player *types.Player, dungeon *types.Dungeon) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

		if callToolResult, err := checkPlayerExists(player); err != nil {
			return callToolResult, err
		}
		// NOTE: generate the ASCII map
		asciiMap := generateASCIIMap(player, dungeon)

		return mcp.NewToolResultText(asciiMap), nil
	}
}

type RoomInfo struct {
	Symbols []string
	Visited bool
}

/*
THE SQUARE DUNGEON OF COMPOSE-AND-DRAGONS
=========================================

    0       1       2       3
  ┌───────┬───────┬───────┬───────┐
3 │ ???   │ ???   │ ???   │ ???   │
  │       │       │       │       │
  │       │       │       │       │
  ├───────┼───────┼───────┼───────┤
2 │       │       │       │ ???   │
  │ [G]   │ [G]   │ [G]   │       │
  │       │       │       │       │
  ├───────┼───────┼───────┼───────┤
1 │       │ ???   │       │ ???   │
  │       │       │ [@][+]│       │
  │       │       │       │       │
  ├───────┼───────┼───────┼───────┤
0 │       │ ???   │ ???   │ ???   │
  │ [E]   │       │       │       │
  │       │       │       │       │
  └───────┴───────┴───────┴───────┘

*/
func generateASCIIMap(player *types.Player, dungeon *types.Dungeon) string {
	var builder strings.Builder

	// STEP 1: Create a grid to track room information
	grid := make([][]RoomInfo, dungeon.Height)
	for i := range grid {
		grid[i] = make([]RoomInfo, dungeon.Width)
		for j := range grid[i] {
			grid[i][j] = RoomInfo{
				Symbols: []string{},
				Visited: false,
			}
		}
	}

	// STEP 2: Fill grid with room data / loop through the dungeon rooms
	for _, room := range dungeon.Rooms {
		x, y := room.Coordinates.X, room.Coordinates.Y
		// Reverse Y to have the correct display (Y=0 at bottom, Y=height-1 at top)
		displayY := dungeon.Height - 1 - y

		roomInfo := RoomInfo{
			Symbols: []string{},
			Visited: room.Visited,
		}

		if room.Visited {
			// Add the symbols in order: Special rooms > Player > Monster > NPC
			if room.IsEntrance {
				roomInfo.Symbols = append(roomInfo.Symbols, "[E]")
			}
			if room.IsExit {
				roomInfo.Symbols = append(roomInfo.Symbols, "[X]")
			}
			if room.HasTreasure {
				roomInfo.Symbols = append(roomInfo.Symbols, "[€]")
			}
			if room.HasMagicPotion {
				roomInfo.Symbols = append(roomInfo.Symbols, "[!]")
			}

			// Player
			if player.Position.X == x && player.Position.Y == y {
				roomInfo.Symbols = append(roomInfo.Symbols, "[@]")
			}

			// Monster
			if room.HasMonster && room.Monster != nil && room.Monster.Kind != "" {
				roomInfo.Symbols = append(roomInfo.Symbols, fmt.Sprintf("[%s]", getMonsterSymbol(room.Monster.Kind)))
			}

			// NPC
			if room.HasNonPlayerCharacter && room.NonPlayerCharacter != nil && room.NonPlayerCharacter.Type != "" {
				roomInfo.Symbols = append(roomInfo.Symbols, fmt.Sprintf("[%s]", getNPCSymbol(room.NonPlayerCharacter.Type)))
			}
		}
		// Update the "information" grid
		grid[displayY][x] = roomInfo
	}

	// STEP 3: Build the ASCII map string
	// Build the title
	title := strings.ToUpper(dungeon.Name)
	builder.WriteString(fmt.Sprintf("%s\n", title))
	builder.WriteString(strings.Repeat("=", len(title)) + "\n\n")

	// Build column headers
	builder.WriteString("    ")
	for x := 0; x < dungeon.Width; x++ {
		builder.WriteString(fmt.Sprintf("%d       ", x))
	}
	builder.WriteString("\n")

	// Build top border
	builder.WriteString("  ┌")
	for x := 0; x < dungeon.Width; x++ {
		builder.WriteString("───────")
		if x < dungeon.Width-1 {
			builder.WriteString("┬")
		}
	}
	builder.WriteString("┐\n")

	// STEP 3: -> STEP 1:
	// BEGIN: Build rows
	for displayY := 0; displayY < dungeon.Height; displayY++ {
		realY := dungeon.Height - 1 - displayY // Convert back to real Y

		// First line of the row (empty or ???)
		builder.WriteString(fmt.Sprintf("%d │", realY))
		for x := 0; x < dungeon.Width; x++ {
			roomInfo := grid[displayY][x]
			if roomInfo.Visited {
				builder.WriteString("       ") // empty space for visited rooms
			} else {
				builder.WriteString(" ???   ") // 7 characters for unvisited rooms & alignment
			}
			if x < dungeon.Width-1 {
				builder.WriteString("│")
			}
		}
		builder.WriteString("│\n")

		// Second line (symbols)
		builder.WriteString("  │")
		for x := 0; x < dungeon.Width; x++ {
			roomInfo := grid[displayY][x]
			if roomInfo.Visited {
				// Build the symbols string with left alignment
				symbolStr := ""
				for _, symbol := range roomInfo.Symbols {
					symbolStr += symbol
				}
				// Align to left in a 7-character space
				builder.WriteString(fmt.Sprintf(" %-6s", symbolStr))
			} else {
				builder.WriteString("       ") // Non visited room
			}
			if x < dungeon.Width-1 {
				builder.WriteString("│")
			}
		}
		builder.WriteString("│\n")

		// Third line (visited marker)
		builder.WriteString("  │")
		for x := 0; x < dungeon.Width; x++ {
			roomInfo := grid[displayY][x]
			if roomInfo.Visited {
				builder.WriteString("  ✓    ")
			} else {
				builder.WriteString("       ")
			}
			if x < dungeon.Width-1 {
				builder.WriteString("│")
			}
		}
		builder.WriteString("│\n")

		// Row separator or bottom border
		if displayY < dungeon.Height-1 {
			builder.WriteString("  ├")
			for x := 0; x < dungeon.Width; x++ {
				builder.WriteString("───────")
				if x < dungeon.Width-1 {
					builder.WriteString("┼")
				}
			}
			builder.WriteString("┤\n")
		}
	} // END: Build rows

	// Bottom border
	builder.WriteString("  └")
	for x := 0; x < dungeon.Width; x++ {
		builder.WriteString("───────")
		if x < dungeon.Width-1 {
			builder.WriteString("┴")
		}
	}
	builder.WriteString("┘\n\n")

	// Add legend
	builder.WriteString("LEGEND:\n=======\n")
	builder.WriteString(fmt.Sprintf("[@] - Player (%s the %s)\n", player.Name, capitalize(player.Class)))
	builder.WriteString("[E] - Entrance\n")

	// STEP 3: -> STEP 2:
	// Add specific NPCs and monsters found
	legendItems := make(map[string]string)
	for _, room := range dungeon.Rooms {
		if room.Visited {
			if room.HasNonPlayerCharacter && room.NonPlayerCharacter != nil && room.NonPlayerCharacter.Type != "" {
				symbol := getNPCSymbol(room.NonPlayerCharacter.Type)
				desc := fmt.Sprintf("%s (%s - %s)", capitalize(string(room.NonPlayerCharacter.Type)), room.NonPlayerCharacter.Name, room.NonPlayerCharacter.Race)
				legendItems[fmt.Sprintf("[%s]", symbol)] = desc
			}
			if room.HasMonster && room.Monster != nil && room.Monster.Kind != "" {
				symbol := getMonsterSymbol(room.Monster.Kind)
				desc := fmt.Sprintf("%s (%s)", capitalize(string(room.Monster.Kind)), room.Monster.Name)
				legendItems[fmt.Sprintf("[%s]", symbol)] = desc
			}
		}
	}

	for symbol, desc := range legendItems {
		builder.WriteString(fmt.Sprintf("%s - %s\n", symbol, desc))
	}

	builder.WriteString(" ✓  - Visited room\n")
	builder.WriteString("??? - Unvisited/Empty room\n\n")

	// STEP 3: -> STEP 3:
	// Room details
	builder.WriteString("ROOM DETAILS:\n=============\n")
	for _, room := range dungeon.Rooms {
		if room.Visited {
			details := fmt.Sprintf("(%d,%d) %s", room.Coordinates.X, room.Coordinates.Y, room.Name)
			if room.IsEntrance {
				details += " - ENTRANCE"
			}
			if room.IsExit {
				details += " - EXIT"
			}
			if room.HasMonster && room.Monster != nil && room.Monster.Kind != "" {
				details += fmt.Sprintf(" - Has %s", capitalize(string(room.Monster.Kind)))
			}
			if room.HasNonPlayerCharacter && room.NonPlayerCharacter != nil && room.NonPlayerCharacter.Type != "" {
				details += fmt.Sprintf(" - Has %s", capitalize(string(room.NonPlayerCharacter.Type)))
			}
			if player.Position.X == room.Coordinates.X && player.Position.Y == room.Coordinates.Y {
				details += " (Current Location)"
			}
			builder.WriteString(details + "\n")
		}
	}
	builder.WriteString("\n")

	// STEP 3: -> STEP 4:
	// Player status
	builder.WriteString("PLAYER STATUS:\n==============\n")
	builder.WriteString(fmt.Sprintf("Name: %s\n", player.Name))
	builder.WriteString(fmt.Sprintf("Class: %s (%s)\n", capitalize(player.Class), capitalize(player.Race)))
	builder.WriteString(fmt.Sprintf("Level: %d\n", player.Level))
	builder.WriteString(fmt.Sprintf("Health: %d/100\n", player.Health))
	builder.WriteString(fmt.Sprintf("Strength: %d\n", player.Strength))
	builder.WriteString(fmt.Sprintf("Experience: %d\n", player.Experience))
	builder.WriteString(fmt.Sprintf("Gold: %d\n\n", player.GoldCoins))

	// STEP 3: -> STEP 5:
	// Current location info
	currentRoom := getCurrentRoom(player, dungeon)
	if currentRoom != nil {
		builder.WriteString(fmt.Sprintf("Current Position: (%d,%d) - %s\n", player.Position.X, player.Position.Y, currentRoom.Name))
		if currentRoom.IsExit {
			builder.WriteString("You have reached the dungeon exit!\n")
		}
	}

	return builder.String()
}

func getMonsterSymbol(kind types.Kind) string {
	switch kind {
	case types.Dragon:
		return "D"
	case types.Troll:
		return "T"
	case types.Orc:
		return "O"
	case types.Goblin:
		return "G"
	case types.Skeleton:
		return "S"
	case types.Zombie:
		return "Z"
	case types.Werewolf:
		return "W"
	case types.Vampire:
		return "V"
	default:
		return "M"
	}
}

func getNPCSymbol(npcType types.NPCType) string {
	switch npcType {
	case types.Merchant:
		return "$"
	case types.Healer:
		return "+"
	case types.Sorcerer:
		return "*"
	case types.Guard:
		return "G"
	default:
		return "N"
	}
}

func getCurrentRoom(player *types.Player, dungeon *types.Dungeon) *types.Room {
	for i := range dungeon.Rooms {
		if dungeon.Rooms[i].Coordinates.X == player.Position.X &&
			dungeon.Rooms[i].Coordinates.Y == player.Position.Y {
			return &dungeon.Rooms[i]
		}
	}
	return nil
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
