package main

import (
	"bufio"
	"context"
	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/helpers"
	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/rag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/compat_oai/openai"
	"github.com/openai/openai-go/option"
)

func main() {
	ctx := context.Background()

	//engineURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	engineURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/v1/")

	similaritySearchLimit := helpers.StringToFloat(helpers.GetEnvOrDefault("SIMILARITY_LIMIT", "0.5"))
	similaritySearchMaxResults := helpers.StringToInt(helpers.GetEnvOrDefault("SIMILARITY_MAX_RESULTS", "2"))

	chunkSize := helpers.StringToInt(helpers.GetEnvOrDefault("CHUNK_SIZE", "512"))
	chunkOverlap := helpers.StringToInt(helpers.GetEnvOrDefault("CHUNK_OVERLAP", "128"))

	temperature := helpers.StringToFloat(helpers.GetEnvOrDefault("SORCERER_MODEL_TEMPERATURE", "0.0"))
	topP := helpers.StringToFloat(helpers.GetEnvOrDefault("SORCERER_MODEL_TOP_P", "0.9"))

	fmt.Printf("🔧 Using model runner at %s\n", engineURL)
	fmt.Printf("🔧 Similarity search limit: %.2f\n", similaritySearchLimit)
	fmt.Printf("🔧 Similarity search max results: %d\n", similaritySearchMaxResults)
	fmt.Printf("🔧 Chunk size: %d\n", chunkSize)
	fmt.Printf("🔧 Chunk overlap: %d\n", chunkOverlap)
	fmt.Printf("🔧 Sorcerer model temperature: %.2f\n", temperature)

	chatModelId := "openai/" + helpers.GetEnvOrDefault("CHAT_MODEL", "ai/qwen2.5:1.5B-F16")
	//embeddingsModelId := "openai/" + helpers.GetEnvOrDefault("EMBEDDING_MODEL", "ai/mxbai-embed-large")
	embeddingsModelId := helpers.GetEnvOrDefault("EMBEDDING_MODEL", "ai/mxbai-embed-large")

	g := genkit.Init(ctx, genkit.WithPlugins(&openai.OpenAI{
		APIKey: "I💙DockerModelRunner",
		Opts: []option.RequestOption{
			option.WithBaseURL(engineURL),
		},
	}))

	messages := []*ai.Message{}

	agentName := helpers.GetEnvOrDefault("SORCERER_NAME", "Elara")

	// ---------------------------------------------------------
	// System Instructions
	// ---------------------------------------------------------
	// ✋ NOTE: load the system instructions from a file
	systemInstructionsPath := helpers.GetEnvOrDefault("SORCERER_SYSTEM_INSTRUCTIONS_PATH", "")
	// Read the content of the file at systemInstructionsContentPath
	systemInstructions, err := helpers.ReadTextFile(systemInstructionsPath)

	if err != nil {
		log.Fatal("😡:", err)
	}

	backgroundAndPersonalityPath := helpers.GetEnvOrDefault("SORCERER_CONTEXT_PATH", "")
	// Make chunk and rag here
	backgroundAndPersonality, err := helpers.ReadTextFile(backgroundAndPersonalityPath)
	if err != nil {
		log.Fatal("😡:", err)
	}
	//chunks := rag.ChunkText(backgroundAndPersonality, chunkSize, chunkOverlap)

	chunks := rag.ChunkWithMarkdownHierarchy(backgroundAndPersonality)

	embedder, vectorStore, err := GenerateEmbeddings(ctx, engineURL, embeddingsModelId, chunks)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Printf("🧠 Created vector store with %d records\n", len(vectorStore.Records))
	// Create the memory vector retriever
	memoryRetriever := rag.DefineMemoryVectorRetriever(g, &vectorStore, embedder)
	fmt.Println("✅ Embeddings generated and vector store ready with", len(vectorStore.Records), "records")

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("🤖🧠 [%s](%s) ask me something - /bye to exit> ", agentName, chatModelId)
		userMessage, _ := reader.ReadString('\n')

		if strings.HasPrefix(userMessage, "/bye") {
			fmt.Println("👋 Bye!")
			break
		}

		if strings.HasPrefix(userMessage, "/history") {
			fmt.Println("📝 Conversation history:")
			for i, msg := range messages {
				// Convert []*ai.Part to string for display
				var parts []string
				for _, part := range msg.Content {
					parts = append(parts, part.Text)
				}
				fmt.Printf("  [%d] %s: %s\n", i, msg.Role, strings.Join(parts, " "))
			}
			continue
		}

		// Retrieve relevant context from the vector store
		similarDocuments, err := RetrieveSimilarDocuments(ctx, userMessage, memoryRetriever)
		if err != nil {
			log.Fatal("😡 when searching for similarities", err)
		}
		//fmt.Println("🧾 Relevant context:\n", similarDocuments)

		messages = append(messages, ai.NewSystemTextMessage(fmt.Sprintf("Relevant context to help you answer the next question:\n%s", similarDocuments)))

		fullResponse, err := genkit.Generate(ctx, g,
			ai.WithModelName(chatModelId),
			ai.WithSystem(systemInstructions),
			// WithMessages sets the messages.
			// These messages will be sandwiched between the system and user prompts.
			ai.WithMessages(
				messages...,
			),
			ai.WithPrompt(userMessage),
			ai.WithConfig(map[string]any{
				"temperature": temperature,
				"top_p":       topP,
			}),

			ai.WithStreaming(func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
				fmt.Print(chunk.Text())
				return nil
			}),
		)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println() // New line after the response
		// Append user message to history
		messages = append(messages, ai.NewUserTextMessage(strings.TrimSpace(userMessage)))
		// Append assistant response to history
		messages = append(messages, ai.NewModelTextMessage(strings.TrimSpace(fullResponse.Text())))

	}

}

func GenerateEmbeddings(ctx context.Context, engineURL string, embeddingModelId string, chunks []string) (ai.Embedder, rag.MemoryVectorStore, error) {
	store := rag.MemoryVectorStore{
		Records: make(map[string]rag.VectorRecord),
	}

	oaiPlugin := &openai.OpenAI{
		APIKey: "I💙DockerModelRunner",
		Opts: []option.RequestOption{
			option.WithBaseURL(engineURL),
		},
	}
	g := genkit.Init(ctx, genkit.WithPlugins(oaiPlugin))
	embedder := oaiPlugin.DefineEmbedder(embeddingModelId, nil)

	for _, chunk := range chunks {
		resp, err := genkit.Embed(ctx, g,
			ai.WithEmbedder(embedder),
			ai.WithTextDocs(chunk),
		)
		if err != nil {
			fmt.Println("😡 Error generating embedding:", err)
			return nil, rag.MemoryVectorStore{}, err
		}
		for i, emb := range resp.Embeddings {
			// Store the embedding in the vector store
			record, errSave := store.Save(rag.VectorRecord{
				Prompt:    chunk,
				Embedding: emb.Embedding,
			})
			if errSave != nil {
				fmt.Println("😡 Error saving vector record:", errSave)
				return nil, rag.MemoryVectorStore{}, errSave
			}
			fmt.Println("💾", i, "Saved record:", record.Prompt, record.Id)

		}
	}
	return embedder, store, nil
	// TODO: save to a JSON file and retrive from there
}

// RetrieveSimilarDocuments performs similarity search and returns relevant context with details
func RetrieveSimilarDocuments(ctx context.Context, query string, retriever ai.Retriever) (string, error) {
	// Create a query document from the user question
	queryDoc := ai.DocumentFromText(query, nil)

	similarityThreshold := helpers.StringToFloat(helpers.GetEnvOrDefault("SIMILARITY_THRESHOLD", "0.5"))
	similarityMaxResults := helpers.StringToInt(helpers.GetEnvOrDefault("SIMILARITY_MAX_RESULTS", "3"))

	// Create a retriever request with custom options
	request := &ai.RetrieverRequest{
		Query: queryDoc,
		Options: rag.MemoryVectorRetrieverOptions{
			Limit:      similarityThreshold,  // Lower similarity threshold to get more results
			MaxResults: similarityMaxResults, // Return top N results
		},
	}

	// Use the memory vector retriever to find similar documents
	retrieveResponse, err := retriever.Retrieve(ctx, request)
	if err != nil {
		return "", err
	}

	similarDocuments := ""
	fmt.Println("--------------------------------------------------")
	fmt.Printf("\n📘 Found %d similar documents:\n", len(retrieveResponse.Documents))
	for i, doc := range retrieveResponse.Documents {
		similarity := doc.Metadata["cosine_similarity"].(float64)
		id := doc.Metadata["id"].(string)
		content := doc.Content[0].Text

		fmt.Printf("%d. ID: %s, Similarity: %.4f\n", i+1, id, similarity)
		fmt.Printf("   Content: %s\n\n", content)

		similarDocuments += content
	}

	fmt.Println("--------------------------------------------------")
	fmt.Println()
	return similarDocuments, nil
}
