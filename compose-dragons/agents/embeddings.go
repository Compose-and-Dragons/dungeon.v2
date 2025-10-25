package agents

import (
	"context"
	"fmt"
	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/msg"
	"github.com/Compose-and-Dragons/dungeon.v2/compose-dragons/rag"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/compat_oai/openai"
	"github.com/openai/openai-go/option"
)

func getEmbedder(ctx context.Context, engineURL string, embeddingModelId string) ai.Embedder {
	oaiPlugin := &openai.OpenAI{
		APIKey: "I💙DockerModelRunner",
		Opts: []option.RequestOption{
			option.WithBaseURL(engineURL),
		},
	}
	genkit.Init(ctx, genkit.WithPlugins(oaiPlugin))
	embedder := oaiPlugin.DefineEmbedder(embeddingModelId, nil)
	return embedder
}

func generateEmbeddings(ctx context.Context, engineURL string, embeddingModelId string, chunks []string) (ai.Embedder, rag.MemoryVectorStore, error) {
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
			msg.DisplayError("😡 Error generating embedding:", err)
			return nil, rag.MemoryVectorStore{}, err
		}
		for i, emb := range resp.Embeddings {
			// Store the embedding in the vector store
			record, errSave := store.Save(rag.VectorRecord{
				Prompt:    chunk,
				Embedding: emb.Embedding,
			})
			if errSave != nil {
				msg.DisplayError("😡 Error saving vector record:", errSave)
				return nil, rag.MemoryVectorStore{}, errSave
			}

			msg.DisplayEmbeddingsMessages(
				fmt.Sprintf("💾 %d Saved record: %s %s", i, record.Prompt, record.Id),
			)
		}
	}
	return embedder, store, nil
	// TODO: save to a JSON file and retrive from there
}
