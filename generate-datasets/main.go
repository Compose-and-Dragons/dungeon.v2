package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"text/template"

	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/agents"
	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/helpers"
	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/rag"
)

type DataSetEntry struct {
	Prompt   string `json:"prompt"`
	Response string `json:"response"`
}

type Settings struct {
	NbDatasetEntries int `json:"nb_dataset_entries"`
	NbIterations     int `json:"nb_iterations"`
}

func main() {
	ctx := context.Background()
	engineURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/v1/")
	chatModelId := "openai/" + helpers.GetEnvOrDefault("CHAT_MODEL", "hf.co/menlo/jan-nano-128k-gguf:q4_k_m")

	fmt.Println("🌍 LLM URL:", engineURL)
	fmt.Println("🤖 Chat Model:", chatModelId)

	agentName := "Dataset Generator"

	temperature := helpers.StringToFloat(helpers.GetEnvOrDefault("MODEL_TEMPERATURE", "0.2"))
	topP := helpers.StringToFloat(helpers.GetEnvOrDefault("MODEL_TOP_P", "0.9"))

	config := agents.Config{
		EngineURL:   engineURL,
		Temperature: temperature,
		TopP:        topP,
		ChatModelId: chatModelId,
	}

	dataSetAgent := agents.NPCAgent{}
	dataSetAgent.Initialize(ctx, config, agentName)

	systemMsg := helpers.GetEnvOrDefault("SYSTEM_INSTRUCTIONS", ``)
	fmt.Println("✅ System Instructions set.", systemMsg)

	contextDirectory := helpers.GetEnvOrDefault("MODEL_CONTEXT_DIRECTORY", "/app/data/chunks")
	allChunks, err := helpers.GetContentFiles(contextDirectory, ".md")
	if err != nil {
		log.Fatal("😡 When reading markdown documents:", err)
	}

	maxGenerationRetries := helpers.StringToInt(helpers.GetEnvOrDefault("MAX_GENERATION_RETRIES", "3"))

	nameOfTheNPC := helpers.GetEnvOrDefault("NPC_NAME", "Elara")

	fmt.Printf("🧙 Generating dataset entries for NPC '%s'...\n", nameOfTheNPC)

	fileIndex := 1

	for _, content := range allChunks {
		//fmt.Println("📄 Loaded chunk:\n", content)
		trainingData := []DataSetEntry{}

		items := rag.SplitTextWithDelimiter(content, "----------")
		context := items[0]
		document := items[1]
		jsonSettings := items[2]

		fmt.Println("🟡 jsonSettings", jsonSettings)

		var settings Settings
		if err := json.Unmarshal([]byte(jsonSettings), &settings); err != nil {
			log.Fatal("😡 Error unmarshaling settings:", err)
			//return nil, errors.New("😡 Error unmarshaling settings")
		}
		promptTemplate := items[3]

		data := struct {
			NameOfTheNPC      string
			NbEntriesPerChunk int
			Chunk             string
		}{
			NameOfTheNPC:      nameOfTheNPC,
			NbEntriesPerChunk: settings.NbDatasetEntries,
			Chunk:             document,
		}

		fmt.Println(strings.Repeat("=", 80))
		fmt.Println("📄 Context:\n", context)
		fmt.Println("📄 Document:\n", document)
		fmt.Printf("⚙️ Settings: %+v\n", settings)
		fmt.Println(strings.Repeat("=", 80))
		userMessage, err := InterpolateString(promptTemplate, data)

		if err != nil {
			log.Fatal("😡 Error interpolating string:", err)
			//return nil, errors.New("😡 Error interpolating string")
		}

		fmt.Println("💬 User Message:\n", userMessage)
		fmt.Println(strings.Repeat("+", 80))

		fmt.Println("Ⓜ️ Generating dataset entries for the iteration 1:")

		entries, jsonString, err := GenerateDatasetEntries(ctx, config, &dataSetAgent, userMessage, nameOfTheNPC, maxGenerationRetries)
		if err != nil {
			log.Println("😡 Error generating dataset entries for chunk, skipping:", err)
			continue
		}

		trainingData = append(trainingData, entries...)

		for i := 1; i < settings.NbIterations; i++ {
			fmt.Printf("Ⓜ️ Generating dataset entries for the iteration %d:\n", i+1)

			userMessage = fmt.Sprintf(`
			Here is the previous response you gave:\n%s\n
			Now, please generate NEW more dataset entries for the same document but with DIFFERENT prompts.:
			%s
			`, jsonString, userMessage)

			entries, newJsonString, err := GenerateDatasetEntries(ctx, config, &dataSetAgent, userMessage, nameOfTheNPC, maxGenerationRetries)
			if err != nil {
				log.Println("😡 Error generating dataset entries for chunk, skipping:", err)
				continue
			}
			jsonString = jsonString + newJsonString
			//jsonString = newJsonString

			trainingData = append(trainingData, entries...)
		}

		dataSetAgent.ResetMessages()

		// Save training data to JSON file
		jsonData, err := json.MarshalIndent(trainingData, "", "  ")
		if err != nil {
			log.Fatal("😡 Error marshaling training data:", err)
		}

		trainDataDirectory := helpers.GetEnvOrDefault("TRAINING_DATA_DIRECTORY", "/app/data")
		trainDataFile := helpers.GetEnvOrDefault("TRAINING_DATA_FILE", "training_data.json")
		if err := helpers.WriteTextFile(trainDataDirectory+"/"+strconv.Itoa(fileIndex)+"."+trainDataFile, string(jsonData)); err != nil {
			log.Fatal("😡 Error writing training data json file:", err)
		}

		fmt.Println("✅ Training data saved:", trainDataDirectory+"/"+strconv.Itoa(fileIndex)+"."+trainDataFile)
		fileIndex++

	}

	// Merge all the generated training data files into one
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🔀 Merging all training data files...")

	trainDataDirectory := helpers.GetEnvOrDefault("TRAINING_DATA_DIRECTORY", "/app/data")
	trainDataFile := helpers.GetEnvOrDefault("TRAINING_DATA_FILE", "training_data.json")

	if err := MergeTrainingDataFiles(trainDataDirectory, trainDataFile); err != nil {
		log.Fatal("😡 Error merging training data files:", err)
	}

}

func GenerateDatasetEntries(ctx context.Context, config agents.Config, dataSetAgent *agents.NPCAgent, userMessage string, nameOfTheNPC string, maxGenerationRetries int) ([]DataSetEntry, string, error) {

	var response string
	maxRetries := maxGenerationRetries

	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {

		response, err = dataSetAgent.JsonCompletion(ctx, config, []DataSetEntry{}, userMessage)
		if err == nil {
			break
		}
		log.Printf("😡 JSON Completion attempt %d/%d failed: %v\n", attempt, maxRetries, err)
	}

	if err != nil {
		log.Printf("😡 JSON Completion failed after %d attempts, skipping iteration\n", maxRetries)
		return nil, "", err
		//continue
	}

	fmt.Println("🗂️ Dataset Entries Response:\n", response)

	var entries []DataSetEntry
	if err := json.Unmarshal([]byte(response), &entries); err != nil {
		log.Fatal("😡 Error unmarshaling dataset entries:", err)
	}

	for _, entry := range entries {
		fmt.Printf("📝 Prompt: %s\n", entry.Prompt)
		fmt.Printf("💡 Response: %s\n", entry.Response)
	}
	return entries, response, nil
}

func process(t *template.Template, vars interface{}) (string, error) {
	var tmplBytes bytes.Buffer

	err := t.Execute(&tmplBytes, vars)
	if err != nil {
		return "", err
	}
	return tmplBytes.String(), nil
}

func InterpolateString(str string, vars interface{}) (string, error) {

	tmpl, err := template.New("tmpl").Parse(str)

	if err != nil {
		return "", err
	}
	output, err := process(tmpl, vars)

	if err != nil {
		return "", err
	}
	return output, nil
}

func MergeTrainingDataFiles(directory, trainDataFile string) error {
	// Find all files ending with trainDataFile
	files, err := helpers.GetAllFilesInDirectory(directory)
	if err != nil {
		return fmt.Errorf("error reading directory: %w", err)
	}

	var mergedData []DataSetEntry

	// Read and merge all matching files
	for _, file := range files {
		if strings.HasSuffix(file, trainDataFile) {
			fmt.Printf("📁 Reading file: %s\n", file)

			content, err := helpers.ReadTextFile(file)
			if err != nil {
				log.Printf("⚠️ Warning: Could not read file %s: %v\n", file, err)
				continue
			}

			var entries []DataSetEntry
			if err := json.Unmarshal([]byte(content), &entries); err != nil {
				log.Printf("⚠️ Warning: Could not unmarshal file %s: %v\n", file, err)
				continue
			}

			mergedData = append(mergedData, entries...)
			fmt.Printf("✅ Added %d entries from %s\n", len(entries), file)
		}
	}

	if len(mergedData) == 0 {
		return fmt.Errorf("no data to merge")
	}

	// Save merged data
	jsonData, err := json.MarshalIndent(mergedData, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling merged data: %w", err)
	}

	outputFile := directory + "/training_" + trainDataFile
	if err := helpers.WriteTextFile(outputFile, string(jsonData)); err != nil {
		return fmt.Errorf("error writing merged file: %w", err)
	}

	fmt.Printf("🎉 Successfully merged %d total entries into: %s\n", len(mergedData), outputFile)
	return nil
}
