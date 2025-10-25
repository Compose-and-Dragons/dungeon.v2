package types

type Player struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
	Class string `json:"class"`
	Race  string `json:"race"`
	Position Coordinates `json:"position"`
	RoomID string `json:"room_id"`
	//Inventory []string `json:"inventory"`
	Health    int      `json:"health"`
	Strength  int      `json:"strength"`
	Experience int      `json:"experience"`
	GoldCoins int      `json:"gold_coins"`
	IsDead   bool     `json:"is_dead"`
}

