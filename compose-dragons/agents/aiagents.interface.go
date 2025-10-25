package agents

import (
	"context"

	"github.com/firebase/genkit/go/ai"
)

type AIAgent interface {
	//GetName() string
	Initialize(ctx context.Context, config Config, name string)
	QuickInitialization()
	SetSystemInstructionsFromFile(systemInstructionsPath string) error
	SetSystemInstructions(systemInstructions string)
	InitializeVectorStoreFromFile(ctx context.Context, config Config, backgroundContextPath string) error
	Completion(ctx context.Context, config Config, userMessage string) (string, error)
	JsonCompletion(ctx context.Context, config Config, outputType any, userMessage string) (string, error)
	SimilaritySearch(ctx context.Context, config Config, userMessage string) (string, error)
	CompletionWithSimilaritySearch(ctx context.Context, config Config, userMessage string) (string, error)
	StreamCompletion(ctx context.Context, config Config, userMessage string, callback ai.ModelStreamCallback) (string, error)
	StreamCompletionWithSimilaritySearch(ctx context.Context, config Config, userMessage string, callback ai.ModelStreamCallback) (string, error)
	DetectAndExecuteToolCalls(ctx context.Context, config Config, userMessage string) (*ToolCallsResult, error)
	DetectAndExecuteToolCallsWithConfirmation(ctx context.Context, config Config, userMessage string) (*ToolCallsResult, error)
	ResetMessages()
	GetHistory() []*ai.Message
	DisplayHistory()
	LoopCompletion(ctx context.Context, config Config)
	DirectExecuteTool(ctx context.Context, config Config, req *ai.ToolRequest) (string, error)
}
