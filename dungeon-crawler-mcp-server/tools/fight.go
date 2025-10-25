package tools

import (
	"context"
	"dungeon-mcp-server/types"
	"fmt"
	"math/rand"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func FightMonsterTool() mcp.Tool {
	return mcp.NewTool("fight_monster",
		mcp.WithDescription(`Fight a monster in your current room using turn-based combat. Each call represents one combat turn with dice rolls for both player and monster.`),
	)
}

func FightMonsterToolHandler(player *types.Player, dungeon *types.Dungeon) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

		// Check if player exists
		if result, err := checkPlayerExists(player); err != nil {
			return result, err
		}

		// Check if player is dead
		if player.IsDead {
			message := "💀 You are dead and cannot fight. You need to be revived first."
			fmt.Println(message)
			return mcp.NewToolResultText(message), fmt.Errorf("player is dead")
		}

		// Find current room
		currentRoom, callToolResult, err := checkPlayerIsInARoom(player, dungeon)
		if err != nil {
			return callToolResult, err
		}

		// Check if there's a monster in the room
		if !currentRoom.HasMonster {
			message := fmt.Sprintf("🏠 There are no monsters to fight in %s.", currentRoom.Name)
			fmt.Println(message)
			return mcp.NewToolResultText(message), nil
		}

		// Get the monster in this room
		monster := currentRoom.Monster
		if monster == nil || monster.IsDead {
			message := fmt.Sprintf("🏠 All monsters in %s are already defeated.", currentRoom.Name)
			fmt.Println(message)
			return mcp.NewToolResultText(message), nil
		}

		// Fight logic/rules
		// COMBAT RULES:
		// 1. Both player and monster roll 2d6 (two six-sided dice)
		// 2. Add strength stat to the dice roll total
		// 3. Higher total wins the combat turn
		// 4. Winner deals damage equal to the difference between totals
		// 5. Combat continues until one combatant reaches 0 health
		// 6. If player wins: gains 10-30 XP and 5-20 gold coins
		// 7. If rolls are tied: no damage is dealt to either combatant

		// Initialize random generator
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		// Combat turn: roll 2d6 for both player and monster
		playerRoll1 := r.Intn(6) + 1
		playerRoll2 := r.Intn(6) + 1
		playerTotal := playerRoll1 + playerRoll2 + player.Strength

		monsterRoll1 := r.Intn(6) + 1
		monsterRoll2 := r.Intn(6) + 1
		monsterTotal := monsterRoll1 + monsterRoll2 + monster.Strength

		message := "⚔️ **COMBAT TURN**\n"
		message += fmt.Sprintf("🎲 %s rolls: %d + %d + %d (strength) = %d\n",
			player.Name, playerRoll1, playerRoll2, player.Strength, playerTotal)
		message += fmt.Sprintf("🎲 %s rolls: %d + %d + %d (strength) = %d\n",
			monster.Name, monsterRoll1, monsterRoll2, monster.Strength, monsterTotal)

		// Determine winner of this turn
		if playerTotal > monsterTotal {
			// Player wins this turn
			damage := playerTotal - monsterTotal
			monster.Health -= damage
			message += fmt.Sprintf("✅ %s wins this turn! %s takes %d damage.\n",
				player.Name, monster.Name, damage)

			if monster.Health <= 0 {
				monster.Health = 0
				monster.IsDead = true

				// Player gains experience and gold
				expGained := 10 + r.Intn(21) // 10-30 experience
				goldGained := 5 + r.Intn(16) // 5-20 gold
				player.Experience += expGained
				player.GoldCoins += goldGained

				// NOTE: you win! 🎉
				message += fmt.Sprintf("💀 %s is defeated!\n", monster.Name)
				message += fmt.Sprintf("⭐ You gain %d experience and %d gold coins!\n", expGained, goldGained)

				// Update room status since monster is dead
				currentRoom.HasMonster = false
			} else {
				message += fmt.Sprintf("❤️ %s has %d health remaining.\n", monster.Name, monster.Health)
			}
		} else if monsterTotal > playerTotal {
			// Monster wins this turn
			damage := monsterTotal - playerTotal
			player.Health -= damage
			message += fmt.Sprintf("💥 %s wins this turn! You take %d damage.\n",
				monster.Name, damage)

			if player.Health <= 0 {
				player.Health = 0
				player.IsDead = true
				// NOTE: you are dead! ☠️
				message += "💀 You have been defeated! You are now dead.\n"
			} else {
				message += fmt.Sprintf("❤️ You have %d health remaining.\n", player.Health)
			}
		} else {
			// Tie - no damage
			message += "⚖️ It's a tie! No damage dealt.\n"
		}

		// Current status
		message += "\n📊 **STATUS:**\n"
		message += fmt.Sprintf("👤 %s: %d health, %d strength\n", player.Name, player.Health, player.Strength)
		if !monster.IsDead {
			message += fmt.Sprintf("👹 %s: %d health, %d strength\n", monster.Name, monster.Health, monster.Strength)
			message += "\n🎯 Call fight_monster again to continue the battle!"
		} else {
			message += fmt.Sprintf("👹 %s: DEFEATED\n", monster.Name)
		}

		fmt.Println(message)
		return mcp.NewToolResultText(message), nil
	}
}
