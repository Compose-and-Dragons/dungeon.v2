package data

// Response schemas for JSON completion mode
// These schemas are used by NPCAgent.JsonCompletion() to enforce structured outputs
//
// Usage example:
//
//	agent.JsonCompletion(ctx, config, Room{}, "Create a dungeon entrance room")
//
// Note: When using tools with conversation history, the Dungeon Master agent
// resets its message history after each interaction to avoid tool call accumulation.
// See dungeon-master/main.go:255 for ResetMessages() implementation.

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

