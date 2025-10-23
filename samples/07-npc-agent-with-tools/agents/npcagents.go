package agents

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"npc-agent-with-tools/helpers"
	"npc-agent-with-tools/rag"
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

	Temperature float64
	TopP        float64

	ChatModelId       string
	EmbeddingsModelId string
	ToolsModelId      string

	Tools []ai.ToolRef
}

// ToolCallsResult holds the result of tool calls detection and execution
type ToolCallsResult struct {
	TotalCalls  int
	Results     []any
	LastMessage string
}

// IMPORTANT: the conversation history is automatically managed
// TODO: add methods to clear history, export history, import history, etc.
type NPCAgent struct {
	Name string

	genKitInstance *genkit.Genkit

	messages []*ai.Message

	systemInstructions      string
	toolsSystemInstructions string
	//backgroundContext  string

	memoryVectorStore rag.MemoryVectorStore
	embedder          ai.Embedder
	memoryRetriever   ai.Retriever
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

func (agent *NPCAgent) SetSystemInstructions(systemInstructions string) {
	agent.systemInstructions = systemInstructions
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
	agent.memoryVectorStore = vectorStore

	memoryRetriever := rag.DefineMemoryVectorRetriever(agent.genKitInstance, &vectorStore, embedder)
	agent.memoryRetriever = memoryRetriever

	fmt.Printf("🧠 Created vector store with %d records\n", len(vectorStore.Records))

	return nil
}

func (agent *NPCAgent) Completion(ctx context.Context, config Config, userMessage string) (string, error) {

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

func (agent *NPCAgent) SimilaritySearch(ctx context.Context, config Config, userMessage string) (string, error) {
	// Retrieve relevant context from the vector store
	similarDocuments, err := retrieveSimilarDocuments(ctx, userMessage, agent.memoryRetriever, config.SimilaritySearchLimit, config.SimilaritySearchMaxResults)
	if err != nil {
		return "", err
	}
	agent.messages = append(agent.messages, ai.NewSystemTextMessage(fmt.Sprintf("Relevant context to help you answer the next question:\n%s", similarDocuments)))

	return similarDocuments, nil
}

func (agent *NPCAgent) CompletionWithSimilaritySearch(ctx context.Context, config Config, userMessage string) (string, error) {

	// Retrieve relevant context from the vector store
	agent.SimilaritySearch(ctx, config, userMessage)

	return agent.Completion(ctx, config, userMessage)

}

func (agent *NPCAgent) StreamCompletion(ctx context.Context, config Config, userMessage string, callback ai.ModelStreamCallback) (string, error) {

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

func (agent *NPCAgent) StreamCompletionWithSimilaritySearch(ctx context.Context, config Config, userMessage string, callback ai.ModelStreamCallback) (string, error) {

	// Retrieve relevant context from the vector store
	agent.SimilaritySearch(ctx, config, userMessage)

	return agent.StreamCompletion(ctx, config, userMessage, callback)

}

func (agent *NPCAgent) DetectAndExecuteToolCalls(ctx context.Context, config Config, userMessage string) (*ToolCallsResult, error) {

	stopped := false
	lastToolAssistantMessage := ""
	totalOfToolsCalls := 0
	toolCallsResults := []any{}

	history := []*ai.Message{}

	fmt.Println("🛠️ Tools index", len(config.Tools), "active tools.")
	for _, t := range config.Tools {
		fmt.Println("   -", t.Name())
	}

	// To avoid repeating the first user message in the history
	// we add it here before entering the loop and using prompt
	history = append(history, ai.NewUserTextMessage(userMessage))


	for !stopped {
		fmt.Printf("\n🔄 Tool detection loop iteration - Current history length: %d\n", len(history))

		resp, err := genkit.Generate(ctx, agent.genKitInstance,
			ai.WithModelName(config.ToolsModelId),
			ai.WithSystem(agent.toolsSystemInstructions),
			ai.WithMessages(history...),
			//ai.WithPrompt(userMessage),
			ai.WithTools(config.Tools...),
			ai.WithToolChoice(ai.ToolChoiceAuto),
			ai.WithReturnToolRequests(true),
		)
		if err != nil {
			fmt.Printf("🔴 [tools] Error: %v\n", err)
		}

		toolRequests := resp.ToolRequests()
		if len(toolRequests) == 0 {
			stopped = true
			lastToolAssistantMessage = resp.Text()
			fmt.Println("✅ No more tool requests, stopping loop")
			break
		}
		fmt.Println("✋ Number of tool requests", len(toolRequests))
		totalOfToolsCalls += len(toolRequests)

		history = append(history, resp.Message)
		fmt.Printf("📥 Added assistant message to history (length now: %d)\n", len(history))

		for _, req := range toolRequests {
			fmt.Println("🛠️ Tool request:", req.Name, "Ref:", req.Ref, "Input:", req.Input)

			// STEP 1: First try to lookup in registered tools (for locally defined tools)
			var tool ai.Tool
			// tool = genkit.LookupTool(agent.genKitInstance, req.Name)
			// if tool != nil {
			// 	fmt.Println("   ✅ Found in genkit registry (local tool)")
			// }

			// STEP 2: If not found, search in config.Tools (for MCP tools)
			//if tool == nil {
			for _, t := range config.Tools {
				if t.Name() == req.Name {
					fmt.Println("   🔍 Found in config.Tools (MCP tool), attempting conversion...")
					// Try to convert ToolRef to Tool
					if toolImpl, ok := t.(ai.Tool); ok {
						tool = toolImpl
						fmt.Println("   ✅ Successfully converted to ai.Tool")
						break
					} else {
						fmt.Println("   ❌ Failed to convert ToolRef to ai.Tool")
					}
				}
			}
			//}

			// STEP 3: If still not found, log error and continue
			if tool == nil {
				fmt.Printf("🔴 tool %q not found\n", req.Name)
				//break // [TODO]: continue?
				continue
			}

			// STEP 4: Ask for tool execution confirmation
			execConfirmation := func() {
				var response string
				for {
					fmt.Printf("Do you want to execute tool %q? (y/n/q): ", req.Name)
					_, err := fmt.Scanln(&response)
					if err != nil {
						fmt.Println("Error reading input:", err)
						continue
					}
					response = strings.ToLower(strings.TrimSpace(response))

					switch response {
					case "q":
						fmt.Println("Exiting the program.")
						stopped = true
						return
					case "y":
						output, err := tool.RunRaw(ctx, req.Input)
						if err != nil {
							log.Fatalf("tool %q execution failed: %v", tool.Name(), err)
						}
						fmt.Println("🤖 Result:", output)

						//toolCallsResults += fmt.Sprintf("Result: %v\n", output)
						toolCallsResults = append(toolCallsResults, map[string]any{
							"tool_name":   req.Name,
							"tool_ref":    req.Ref,
							"tool_output": output,
						})

						part := ai.NewToolResponsePart(&ai.ToolResponse{
							Name:   req.Name,
							Ref:    req.Ref,
							Output: output,
						})
						//fmt.Println("✅", output)
						history = append(history, ai.NewMessage(ai.RoleTool, nil, part))
						fmt.Printf("   📜 History length now: %d\n", len(history))

						return
					case "n":
						fmt.Println("⏩ Skipping tool execution.", req.Name, req.Ref)

						//toolCallsResults += fmt.Sprintf("Result: tool %v execution cancelled by user\n", req.Name)
						toolCallsResults = append(toolCallsResults, map[string]any{
							"tool_name":   req.Name,
							"tool_ref":    req.Ref,
							"tool_output": "Tool execution cancelled by user",
						})

						// Add tool response indicating the tool was not executed
						part := ai.NewToolResponsePart(&ai.ToolResponse{
							Name:   req.Name,
							Ref:    req.Ref,
							Output: map[string]any{"error": "Tool execution cancelled by user"},
						})
						history = append(history, ai.NewMessage(ai.RoleTool, nil, part))
						fmt.Printf("   📜 History length now: %d\n", len(history))

						return
					default:
						fmt.Println("Please enter 'y' or 'n'.")
						continue
					}

				}

			}
			execConfirmation()

			fmt.Println(strings.Repeat("-", 20))
			fmt.Println("📜 Tools History now has", len(history), "messages")
			fmt.Println(strings.Repeat("-", 20))
		}

	}
	//fmt.Println("🎉 Final response:\n", lastToolAssistantMessage)

	return &ToolCallsResult{
		TotalCalls:  totalOfToolsCalls,
		Results:     toolCallsResults,
		LastMessage: lastToolAssistantMessage,
	}, nil
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

		agent.StreamCompletionWithSimilaritySearch(ctx, config, userMessage, func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
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
