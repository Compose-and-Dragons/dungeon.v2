package types


type Kind string

const (
	Skeleton Kind = "skeleton"
	Zombie   Kind = "zombie"
	Goblin   Kind = "goblin"
	Orc      Kind = "orc"
	Troll    Kind = "troll"
	Dragon   Kind = "dragon"
	Werewolf Kind = "werewolf"
	Vampire  Kind = "vampire"

	Nothing Kind = "nothing"
)

type Monster struct {
	Kind       Kind   `json:"kind"`
	Name       string `json:"name"`
	Description string `json:"description"`
	Health     int    `json:"health"`
	Strength   int    `json:"strength"`
	Position Coordinates `json:"position"`
	RoomID   string      `json:"room_id"`
	// RoomID     string `json:"room_id"`
	// Position   Coordinates `json:"position"`
	// Experience int    `json:"experience"`
	// GoldCoins  int    `json:"gold_coins"`
	IsDead     bool   `json:"is_dead"`
}