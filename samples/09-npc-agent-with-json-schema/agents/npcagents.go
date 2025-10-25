package agents

import (
	"bufio"
	"context"
	"fmt"
	"npc-agent-with-json-schema/helpers"
	"npc-agent-with-json-schema/msg"
	"npc-agent-with-json-schema/rag"
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
		APIKey: "I💙DockerModelRunner",
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

	msg.DisplayEmbeddingsMessages(
		fmt.Sprintf("🧠 Created vector store with %d records\n", len(vectorStore.Records)),
	)

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

func (agent *NPCAgent) JsonCompletion(ctx context.Context, config Config, outputType any, userMessage string) (string, error) {
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
		ai.WithOutputType(outputType),
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

// executeTool executes a tool without confirmation
//
// Flow:
//
//	┌─────────────────┐
//	│  executeTool    │
//	└────────┬────────┘
//	         │
//	         ▼
//	┌─────────────────┐
//	│ tool.RunRaw()   │
//	└────────┬────────┘
//	         │
//	         ▼
//	┌─────────────────┐
//	│ Append result   │
//	│ to history      │
//	└─────────────────┘
func (agent *NPCAgent) executeTool(ctx context.Context, req *ai.ToolRequest, tool ai.Tool, history *[]*ai.Message, toolCallsResults *[]any) error {
	output, err := tool.RunRaw(ctx, req.Input)
	if err != nil {
		msg.DisplayError(fmt.Sprintf("😡 tool %q execution failed:", tool.Name()), err)
		return err
	}

	msg.DisplayToolMessages(fmt.Sprintf("🤖 Result: %v", output))

	*toolCallsResults = append(*toolCallsResults, map[string]any{
		"tool_name":   req.Name,
		"tool_ref":    req.Ref,
		"tool_output": output,
	})

	// Add tool response to history
	part := ai.NewToolResponsePart(&ai.ToolResponse{
		Name:   req.Name,
		Ref:    req.Ref,
		Output: output,
	})

	// Append tool response to history
	*history = append(*history, ai.NewMessage(ai.RoleTool, nil, part))
	return nil
}

// executeToolWithConfirmation asks for user confirmation before executing a tool
//
// Flow:
//
//	┌─────────────────────────-─┐
//	│executeToolWithConfirmation│
//	└────────────┬────────────-─┘
//	             │
//	             ▼
//	┌────────────────────────-──┐
//	│ Prompt user: y/n/q?       │
//	└────────────┬────────────-─┘
//	             │
//	    ┌────────┼────────┐
//	    │        │        │
//	    ▼        ▼        ▼
//	 ┌───┐   ┌───┐   ┌───┐
//	 │ y │   │ n │   │ q │
//	 └─┬─┘   └─┬─┘   └─┬─┘
//	   │       │       │
//	   │       │       └──► Set stopped=true
//	   │       │
//	   │       └──► Append "cancelled" to history
//	   │
//	   └──► Call executeTool()
func (agent *NPCAgent) executeToolWithConfirmation(ctx context.Context, req *ai.ToolRequest, tool ai.Tool, history *[]*ai.Message, toolCallsResults *[]any, stopped *bool) {
	var response string
	for {
		fmt.Printf("Do you want to execute tool %q? (y/n/q): ", req.Name)
		_, err := fmt.Scanln(&response)
		if err != nil {
			msg.DisplayError("😡 Error reading input:", err)
			continue
		}
		response = strings.ToLower(strings.TrimSpace(response))

		switch response {
		case "q":
			fmt.Println("👋 Exiting the program.")
			*stopped = true
			return
		case "y":
			agent.executeTool(ctx, req, tool, history, toolCallsResults)
			return
		case "n":

			fmt.Println("⏩ Skipping tool execution.", req.Name, req.Ref)

			*toolCallsResults = append(*toolCallsResults, map[string]any{
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

			*history = append(*history, ai.NewMessage(ai.RoleTool, nil, part))
			return
		default:
			fmt.Println("Please enter 'y' or 'n'.")
			continue
		}

	}
}

// toolExecutorFunc is a function type for executing tools
type toolExecutorFunc func(ctx context.Context, req *ai.ToolRequest, tool ai.Tool, history *[]*ai.Message, toolCallsResults *[]any, stopped *bool)

// detectAndExecuteToolCallsLoop is the core loop for detecting and executing tool calls
// It takes an executor function as parameter to customize the execution behavior
//
// Flow:
//
//	┌────────────────────────────────┐
//	│ detectAndExecuteToolCallsLoop  │
//	└──────────────┬─────────────────┘
//	               │
//	               ▼
//	┌──────────────────────────────┐
//	│ Add user message to history  │
//	└──────────────┬───────────────┘
//	               │
//	       ┌───────▼───────┐
//	       │  Loop: while  │
//	       │  !stopped     │
//	       └───────┬───────┘
//	               │
//	               ▼
//	┌──────────────────────────────┐
//	│ Generate AI response with    │
//	│ tool detection               │
//	└──────────────┬───────────────┘
//	               │
//	       ┌───────▼───────────┐
//	       │ Tool requests     │
//	       │ found?            │
//	       └───────┬───────────┘
//	               │
//	       ┌───────┼───────┐
//	       │               │
//	       ▼ No            ▼ Yes
//	┌──────────┐    ┌──────────────────┐
//	│  Stop    │    │ For each tool:   │
//	│  loop    │    │ 1. Find tool     │
//	└──────────┘    │ 2. Call executor │
//	                └────────┬─────────┘
//	                         │
//	                         └──► Loop back
func (agent *NPCAgent) detectAndExecuteToolCallsLoop(ctx context.Context, config Config, userMessage string, executor toolExecutorFunc) (*ToolCallsResult, error) {

	stopped := false
	lastToolAssistantMessage := ""
	totalOfToolsCalls := 0
	toolCallsResults := []any{}

	history := []*ai.Message{}

	// Only displayed if enabled via env var ...
	displayToolsList(config.Tools)

	// IMPORTANT:
	// To avoid repeating the first user message in the history
	// we add it here before entering the loop and using prompt
	history = append(history, ai.NewUserTextMessage(userMessage))

	for !stopped {
		//msg.DisplayToolMessages(fmt.Sprintf("\n🔄 Tool detection loop iteration - Current history length: %d\n", len(history)))

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
			msg.DisplayError("🔴 [tools] Error:", err)
			// break the loop on error
			stopped = true
			break
		}

		toolRequests := resp.ToolRequests()
		// If no tool requests, we can stop the loop
		if len(toolRequests) == 0 {
			stopped = true
			lastToolAssistantMessage = resp.Text()

			msg.DisplayToolMessages("✅ No more tool requests, stopping loop")
			break
		}
		//msg.DisplayToolMessages(fmt.Sprintf("✋ Number of tool requests: %v", len(toolRequests)))

		totalOfToolsCalls += len(toolRequests)

		// Append the assistant message with tool requests to history
		history = append(history, resp.Message)

		// BEGIN: [TOOL CALLS] detection loop
		for _, req := range toolRequests {
			// STEP 1: find the tool by name
			msg.DisplayToolMessages(fmt.Sprintf("🛠️ Tool request: %s Args: %v", req.Name, req.Input))
			var tool ai.Tool
			// tool = genkit.LookupTool(agent.genKitInstance, req.Name)

			for _, t := range config.Tools {
				if t.Name() == req.Name {
					// Try to convert ToolRef to Tool
					if toolImpl, ok := t.(ai.Tool); ok {
						tool = toolImpl
						// ✅ Successfully converted to ai.Tool"
						break
					}
					// else: ❌ Failed to convert ToolRef to ai.Tool")
				}
			}

			// If not found, log error and continue
			if tool == nil {
				msg.DisplayToolMessages(fmt.Sprintf("🔴 tool %q not found\n", req.Name))
				//break // [TODO]: continue?
				continue
			}

			// STEP 2: Execute tool using the provided executor
			executor(ctx, req, tool, &history, &toolCallsResults, &stopped)
		}

	} // END: of [TOOL CALLS] detection loop

	// [TOOL CALL RESULT]
	return &ToolCallsResult{
		TotalCalls:  totalOfToolsCalls,
		Results:     toolCallsResults,
		LastMessage: lastToolAssistantMessage,
	}, nil
}

// DetectAndExecuteToolCalls detects and executes tool calls automatically (no confirmation)
//
// Flow:
//
//	┌──────────────────────────────┐
//	│ DetectAndExecuteToolCalls    │
//	└──────────────┬───────────────┘
//	               │
//	               │ Creates executor lambda
//	               │ that calls executeTool
//	               │
//	               ▼
//	┌──────────────────────────────┐
//	│detectAndExecuteToolCallsLoop │
//	└──────────────┬───────────────┘
//	               │
//	               │ For each tool request,
//	               │ calls executor
//	               │
//	               ▼
//	┌──────────────────────────────┐
//	│ executeTool                  │
//	│ (automatic execution)        │
//	└──────────────┬───────────────┘
//	               │
//	               ▼
//	┌──────────────────────────────┐
//	│ tool.RunRaw()                │
//	│ Append result to history     │
//	└──────────────────────────────┘
func (agent *NPCAgent) DetectAndExecuteToolCalls(ctx context.Context, config Config, userMessage string) (*ToolCallsResult, error) {
	executor := func(ctx context.Context, req *ai.ToolRequest, tool ai.Tool, history *[]*ai.Message, toolCallsResults *[]any, stopped *bool) {
		agent.executeTool(ctx, req, tool, history, toolCallsResults)
	}
	return agent.detectAndExecuteToolCallsLoop(ctx, config, userMessage, executor)
}

// DetectAndExecuteToolCallsWithConfirmation detects and executes tool calls with user confirmation
//
// Flow:
//
//	┌──────────────────────────────┐
//	│DetectAndExecuteToolCalls     │
//	│WithConfirmation              │
//	└──────────────┬───────────────┘
//	               │
//	               │ Passes executeToolWithConfirmation
//	               │ as executor
//	               │
//	               ▼
//	┌──────────────────────────────┐
//	│detectAndExecuteToolCallsLoop │
//	└──────────────┬───────────────┘
//	               │
//	               │ For each tool request,
//	               │ calls executor
//	               │
//	               ▼
//	┌──────────────────────────────┐
//	│executeToolWithConfirmation   │
//	└──────────────┬───────────────┘
//	               │
//	               │ Ask user: y/n/q?
//	               │
//	    ┌──────────┼──────────┐
//	    │          │          │
//	    ▼          ▼          ▼
//	 ┌───┐      ┌───┐      ┌───┐
//	 │ y │      │ n │      │ q │
//	 └─┬─┘      └─┬─┘      └─┬─┘
//	   │          │          │
//	   │          │          └──► Set stopped=true
//	   │          │
//	   │          └──► Append "cancelled" to history
//	   │
//	   └──► Call executeTool()
func (agent *NPCAgent) DetectAndExecuteToolCallsWithConfirmation(ctx context.Context, config Config, userMessage string) (*ToolCallsResult, error) {
	return agent.detectAndExecuteToolCallsLoop(ctx, config, userMessage, agent.executeToolWithConfirmation)
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
