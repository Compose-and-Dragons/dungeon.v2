package agents

import (
	"bufio"
	"context"
	"fmt"
	"npc-agent/helpers"
	"npc-agent/rag"
	"os"
	"strings"

	"github.com/firebase/genkit/go/ai"

	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/compat_oai/openai"
	"github.com/openai/openai-go/option"
)

type Config struct {
	EngineURL                  string
	SimilaritySearchLimit      float64
	SimilaritySearchMaxResults int

	/*
		similarityThreshold := helpers.StringToFloat(helpers.GetEnvOrDefault("SIMILARITY_THRESHOLD", "0.5"))
		similarityMaxResults := helpers.StringToInt(helpers.GetEnvOrDefault("SIMILARITY_MAX_RESULTS", "3"))

	*/

	Temperature float64
	TopP        float64

	ChatModelId       string
	EmbeddingsModelId string
}

type NPCAgent struct {
	Name string

	genKitInstance *genkit.Genkit

	messages []*ai.Message

	systemInstructions string
	//backgroundContext  string

	vectorStore     rag.MemoryVectorStore
	embedder        ai.Embedder
	memoryRetriever ai.Retriever
}

func (agent *NPCAgent) Initialize(ctx context.Context, config Config, name string) {
	// Initialization logic for the NPC agent
	oaiPlugin := &openai.OpenAI{
		APIKey: "tada",
		Opts: []option.RequestOption{
			option.WithBaseURL(config.EngineURL),
		},
	}
	g := genkit.Init(ctx, genkit.WithPlugins(oaiPlugin))
	agent.genKitInstance = g

	agent.Name = name
	agent.messages = []*ai.Message{}

}

func (agent *NPCAgent) QuickInitialization() {
	// Initialization logic for the NPC agent
	// TODO: implement if needed
}

func (agent *NPCAgent) SetSystemInstructionsFromFile(systemInstructionsPath string) error {
	systemInstructions, err := helpers.ReadTextFile(systemInstructionsPath)
	if err != nil {
		return err
	}
	agent.systemInstructions = systemInstructions
	return nil
}

func (agent *NPCAgent) InitializeVectorStoreFromFile(ctx context.Context, config Config, backgroundContextPath string) error {
	backgroundContext, err := helpers.ReadTextFile(backgroundContextPath)
	if err != nil {
		return err
	}
	chunks := rag.ChunkWithMarkdownHierarchy(backgroundContext)

	embedder, vectorStore, err := generateEmbeddings(ctx, config.EngineURL, config.EmbeddingsModelId, chunks)
	if err != nil {
		return err
	}
	agent.embedder = embedder
	agent.vectorStore = vectorStore

	memoryRetriever := rag.DefineMemoryVectorRetriever(agent.genKitInstance, &vectorStore, embedder)
	agent.memoryRetriever = memoryRetriever

	fmt.Printf("🧠 Created vector store with %d records\n", len(vectorStore.Records))

	return nil
}

func (agent *NPCAgent) StreamCompletion(ctx context.Context, config Config, userMessage string, callback ai.ModelStreamCallback) (string, error) {

	// Retrieve relevant context from the vector store
	similarDocuments, err := retrieveSimilarDocuments(ctx, userMessage, agent.memoryRetriever, config.SimilaritySearchLimit, config.SimilaritySearchMaxResults)
	if err != nil {
		return "", err
	}

	agent.messages = append(agent.messages, ai.NewSystemTextMessage(fmt.Sprintf("Relevant context to help you answer the next question:\n%s", similarDocuments)))

	fullResponse, err := genkit.Generate(ctx, agent.genKitInstance,
		ai.WithModelName(config.ChatModelId),
		ai.WithSystem(agent.systemInstructions),
		// WithMessages sets the messages.
		// These messages will be sandwiched between the system and user prompts.
		ai.WithMessages(
			agent.messages...,
		),
		ai.WithPrompt(userMessage),
		ai.WithConfig(map[string]any{
			"temperature": config.Temperature,
			"top_p":       config.TopP,
		}),

		ai.WithStreaming(callback),
	)

	if err != nil {
		return "", err
	}

	// Append user message to history
	agent.messages = append(agent.messages, ai.NewUserTextMessage(strings.TrimSpace(userMessage)))
	// Append assistant response to history
	agent.messages = append(agent.messages, ai.NewModelTextMessage(strings.TrimSpace(fullResponse.Text())))

	return fullResponse.Text(), nil
}

func (agent *NPCAgent) LoopCompletion(ctx context.Context, config Config) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("🤖🧠 [%s](%s) ask me something - /bye to exit> ", agent.Name, config.ChatModelId)
		userMessage, _ := reader.ReadString('\n')

		if strings.HasPrefix(userMessage, "/bye") {
			fmt.Println("👋 Bye!")
			break
		}

		if strings.HasPrefix(userMessage, "/history") {
			fmt.Println("📝 Conversation history:")
			for i, msg := range agent.messages {
				// Convert []*ai.Part to string for display
				var parts []string
				for _, part := range msg.Content {
					parts = append(parts, part.Text)
				}
				fmt.Printf("  [%d] %s: %s\n", i, msg.Role, strings.Join(parts, " "))
			}
			continue
		}

		agent.StreamCompletion(ctx, config, userMessage, func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
			fmt.Print(chunk.Text())
			return nil
		})

		fmt.Println()
		fmt.Println()

	}
}

func generateEmbeddings(ctx context.Context, engineURL string, embeddingModelId string, chunks []string) (ai.Embedder, rag.MemoryVectorStore, error) {
	store := rag.MemoryVectorStore{
		Records: make(map[string]rag.VectorRecord),
	}

	oaiPlugin := &openai.OpenAI{
		APIKey: "tada",
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

func retrieveSimilarDocuments(ctx context.Context, query string, retriever ai.Retriever, similarityThreshold float64, similarityMaxResults int) (string, error) {
	// Create a query document from the user question
	queryDoc := ai.DocumentFromText(query, nil)

	//similarityThreshold := helpers.StringToFloat(helpers.GetEnvOrDefault("SIMILARITY_THRESHOLD", "0.5"))
	//similarityMaxResults := helpers.StringToInt(helpers.GetEnvOrDefault("SIMILARITY_MAX_RESULTS", "3"))

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
	// TODO: make a log helper
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
