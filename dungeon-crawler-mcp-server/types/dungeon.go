package types

type Dungeon struct {
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	Width          int         `json:"width"`
	Height         int         `json:"height"`
	Rooms          []Room      `json:"rooms"`
	EntranceCoords Coordinates `json:"entrance_coords"`
	ExitCoords     Coordinates `json:"exit_coords"`
	 
}

type Coordinates struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Room struct {
	ID                    string              `json:"id"`
	Name                  string              `json:"name"`
	Description           string              `json:"description"`
	IsEntrance            bool                `json:"is_entrance"`
	IsExit                bool                `json:"is_exit"`
	Coordinates           Coordinates         `json:"coordinates"`
	Visited               bool                `json:"visited"`
	HasMonster            bool                `json:"has_monster"`
	Monster               *Monster            `json:"monster,omitempty"`
	HasNonPlayerCharacter bool                `json:"has_non_player_character"`
	HasTreasure           bool                `json:"has_treasure"`
	GoldCoins             int                 `json:"gold_coins"`
	HasMagicPotion        bool                `json:"has_magic_potion"`
	RegenerationHealth    int                 `json:"regeneration_health"`
	NonPlayerCharacter    *NonPlayerCharacter `json:"non_player_character,omitempty"`
	//IsThePlayerHere       bool                `json:"is_the_player_here"`
}
