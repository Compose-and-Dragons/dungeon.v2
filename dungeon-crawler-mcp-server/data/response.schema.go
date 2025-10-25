package data


type Room struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Monster struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Health      int    `json:"health"`
	Strength    int    `json:"strength"`
	Kind        string `json:"kind"`
}

