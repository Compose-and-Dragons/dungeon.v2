package main

import (
	"context"
	"fmt"
	"log"
	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/helpers"
	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/agents"
	"strings"
	"time"

	"math/rand"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/compat_oai/openai"
)

type DiceRollInput struct {
	NumDice  int `json:"num_dice"`
	NumFaces int `json:"num_faces"`
}

type DiceRollResult struct {
	Rolls []int `json:"rolls"`
	Total int   `json:"total"`
}

type CharacterNameInput struct {
	Race string `json:"race"`
}

type CharacterNameResult struct {
	Name string `json:"name"`
	Race string `json:"race"`
}

func main() {
	ctx := context.Background()

	//engineURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	engineURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/v1/")
	// IMPORTANT: prefix with "openai/" to use the OpenAI plugin TODO: make this automatic
	chatModelId := "openai/" + helpers.GetEnvOrDefault("CHAT_MODEL", "ai/qwen2.5:1.5B-F16")
	embeddingsModelId := helpers.GetEnvOrDefault("EMBEDDING_MODEL", "ai/mxbai-embed-large")
	toolsModelId := "openai/" + helpers.GetEnvOrDefault("TOOLS_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")

	agentName := helpers.GetEnvOrDefault("SORCERER_NAME", "Elara")

	similaritySearchLimit := helpers.StringToFloat(helpers.GetEnvOrDefault("SIMILARITY_LIMIT", "0.5"))
	similaritySearchMaxResults := helpers.StringToInt(helpers.GetEnvOrDefault("SIMILARITY_MAX_RESULTS", "2"))

	temperature := helpers.StringToFloat(helpers.GetEnvOrDefault("SORCERER_MODEL_TEMPERATURE", "0.0"))
	topP := helpers.StringToFloat(helpers.GetEnvOrDefault("SORCERER_MODEL_TOP_P", "0.9"))

	g := genkit.Init(ctx, genkit.WithPlugins(&openai.OpenAI{
		APIKey: "I💙DockerModelRunner",
		// Opts: []option.RequestOption{
		// 	option.WithBaseURL("http://localhost:12434/engines/v1/"),
		// },
	}))

	// Define tools
	diceRollTool := genkit.DefineTool(g, "roll_dice", "Roll n dice with n faces each",
		func(ctx *ai.ToolContext, input DiceRollInput) (DiceRollResult, error) {
			return rollDice(input.NumDice, input.NumFaces), nil
		},
	)

	characterNameTool := genkit.DefineTool(g, "generate_character_name", "Generate a D&D character name for a specific race",
		func(ctx *ai.ToolContext, input CharacterNameInput) (CharacterNameResult, error) {
			return generateCharacterName(input.Race), nil
		},
	)

	config := agents.Config{
		EngineURL:                  engineURL,
		SimilaritySearchLimit:      similaritySearchLimit,
		SimilaritySearchMaxResults: similaritySearchMaxResults,
		Temperature:                temperature,
		TopP:                       topP,
		ChatModelId:                chatModelId,
		EmbeddingsModelId:          embeddingsModelId,
		ToolsModelId:               toolsModelId,
		Tools:                      []ai.ToolRef{diceRollTool, characterNameTool},
	}

	sorcererAgent := agents.NPCAgent{}
	sorcererAgent.Initialize(ctx, config, agentName)

	// ---------------------------------------------------------
	// System Instructions
	// ---------------------------------------------------------
	// systemInstructionsPath := helpers.GetEnvOrDefault("SORCERER_SYSTEM_INSTRUCTIONS_PATH", "")
	// err := sorcererAgent.SetSystemInstructionsFromFile(systemInstructionsPath)
	// if err != nil {
	// 	log.Fatal("😡:", err)
	// }

	// Create system message
	systemMsg := `
	You are a helpful D&D assistant that can roll dice and generate character names.
	Use the appropriate tools when asked to roll dice or generate character names.
	`
	sorcererAgent.SetSystemInstructions(systemMsg)

	// backgroundAndPersonalityPath := helpers.GetEnvOrDefault("SORCERER_CONTEXT_PATH", "")
	// err := sorcererAgent.InitializeVectorStoreFromFile(ctx, config, backgroundAndPersonalityPath)
	// if err != nil {
	// 	log.Fatal("😡:", err)
	// }

	toolCallsResult, err := sorcererAgent.DetectAndExecuteToolCallsWithConfirmation(ctx, config, `
		Roll 3 dices with 6 faces each. 
		Then generate a character name for an elf.
		Finally, roll 2 dices with 8 faces each.
		After that, generate a character name for a dwarf.
	`)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println("🛠️ Total calls:", toolCallsResult.TotalCalls)
	fmt.Println("🛠️ Results:\n", toolCallsResult.Results)
	fmt.Println("🛠️ Final Answer:\n", toolCallsResult.LastMessage)

	fmt.Println(strings.Repeat("=", 50))

	toolCallsResult, err = sorcererAgent.DetectAndExecuteToolCalls(ctx, config, `
		Generate a character name for a human.
		Finally, roll 2 dices with 10 faces each.
	`)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println("🛠️ Total calls:", toolCallsResult.TotalCalls)
	fmt.Println("🛠️ Results:\n", toolCallsResult.Results)
	fmt.Println("🛠️ Final Answer:\n", toolCallsResult.LastMessage)
	//sorcererAgent.LoopCompletion(ctx, config)

}

func rollDice(numDice, numFaces int) DiceRollResult {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rolls := make([]int, numDice)
	total := 0

	for i := 0; i < numDice; i++ {
		roll := r.Intn(numFaces) + 1
		rolls[i] = roll
		total += roll
	}

	return DiceRollResult{
		Rolls: rolls,
		Total: total,
	}
}

func generateCharacterName(race string) CharacterNameResult {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	namesByRace := map[string][]string{
		"elf":      {"Aerdrie", "Ahvonna", "Aramil", "Aranea", "Berrian", "Caelynn", "Carric", "Dayereth", "Enna", "Galinndan"},
		"dwarf":    {"Adrik", "Baern", "Darrak", "Eberk", "Fargrim", "Gardain", "Harbek", "Kildrak", "Morgran", "Thorek"},
		"human":    {"Aerdrie", "Aramil", "Berris", "Cithreth", "Dayereth", "Enna", "Galinndan", "Hadarai", "Immeral", "Lamlis"},
		"halfling": {"Alton", "Ander", "Bernie", "Bobbin", "Cade", "Callus", "Corrin", "Dannad", "Garret", "Lindal"},
		"orc":      {"Gash", "Gell", "Henk", "Holg", "Imsh", "Keth", "Krusk", "Mhurren", "Ront", "Shump"},
		"tiefling": {"Akmenos", "Amnon", "Barakas", "Damakos", "Ekemon", "Iados", "Kairon", "Leucis", "Melech", "Mordai"},
	}

	raceLower := strings.ToLower(race)
	names, exists := namesByRace[raceLower]
	if !exists {
		names = namesByRace["human"] // Default to human names
	}

	selectedName := names[r.Intn(len(names))]

	return CharacterNameResult{
		Name: selectedName,
		Race: race,
	}
}
