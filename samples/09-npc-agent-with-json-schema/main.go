package main

import (
	"context"
	"fmt"
	"npc-agent-with-json-schema/agents"
	"npc-agent-with-json-schema/helpers"
)

func main() {
	ctx := context.Background()

	//engineURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	engineURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/v1/")
	// IMPORTANT: prefix with "openai/" to use the OpenAI plugin TODO: make this automatic
	chatModelId := "openai/" + helpers.GetEnvOrDefault("CHAT_MODEL", "ai/qwen2.5:1.5B-F16")

	fmt.Println("🌍 LLM URL:", engineURL)
	fmt.Println("🤖 Chat Model:", chatModelId)

	agentName := helpers.GetEnvOrDefault("SORCERER_NAME", "Elara")

	temperature := helpers.StringToFloat(helpers.GetEnvOrDefault("SORCERER_MODEL_TEMPERATURE", "0.0"))
	topP := helpers.StringToFloat(helpers.GetEnvOrDefault("SORCERER_MODEL_TOP_P", "0.9"))

	config := agents.Config{
		EngineURL:   engineURL,
		Temperature: temperature,
		TopP:        topP,
		ChatModelId: chatModelId,
	}

	sorcererAgent := agents.NPCAgent{}
	sorcererAgent.Initialize(ctx, config, agentName)

	// Create system message
	systemMsg := `
	You are a helpful D&D assistant that can roll dice and generate character names.
	Use the appropriate tools when asked to roll dice or generate character names.
	`
	sorcererAgent.SetSystemInstructions(systemMsg)

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

	response, err := sorcererAgent.JsonCompletion(ctx, config, Room{}, `
		Create an dungeon entrance room with a name and a short description.
	`)
	if err != nil {
		panic(err)
	}
	fmt.Println("🏰 Room Response:\n", response)

	response, err = sorcererAgent.JsonCompletion(ctx, config, Monster{}, `
		Create a new monster with a name and a short description.
		The value of health should be between 50 and 200.
		The value of strength should be between 10 and 50.
		The value of kind should be one of: skeleton, zombie, goblin, orc, troll, dragon, werewolf, vampire.
	`)
	if err != nil {
		panic(err)
	}
	fmt.Println("👹 Monster Response:\n", response)	

		response, err = sorcererAgent.JsonCompletion(ctx, config, Monster{}, `
		Create a new monster with a name and a short description..
	`)
	if err != nil {
		panic(err)
	}
	fmt.Println("🐺 Monster Response:\n", response)

}
