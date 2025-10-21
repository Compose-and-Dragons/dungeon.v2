package rag

import (
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// MemoryVectorRetrieverOptions defines the options for the memory vector retriever
type MemoryVectorRetrieverOptions struct {
	Limit      float64 // Minimum similarity threshold
	MaxResults int     // Maximum number of results to return
}

// DefineMemoryVectorRetriever creates a memory vector retriever using the MemoryVectorStore
func DefineMemoryVectorRetriever(g *genkit.Genkit, vectorStore *MemoryVectorStore, embedder ai.Embedder) ai.Retriever {
	return genkit.DefineRetriever(
		g, "memoryVectorStoreRetriever", nil,
		func(ctx context.Context, req *ai.RetrieverRequest) (*ai.RetrieverResponse, error) {
			// Handle memory vector retriever options
			opts, ok := req.Options.(MemoryVectorRetrieverOptions)
			if !ok {
				// Set default values if no options provided
				opts = MemoryVectorRetrieverOptions{
					Limit:      0.5, // Default similarity threshold
					MaxResults: 5,   // Default max results
				}
			}

			// Set default values for zero values
			if opts.Limit == 0 {
				opts.Limit = 0.5
			}
			if opts.MaxResults == 0 {
				opts.MaxResults = 5
			}

			// Extract query text from the document
			var queryText string
			if req.Query != nil && len(req.Query.Content) > 0 {
				queryText = req.Query.Content[0].Text
			} else {
				return &ai.RetrieverResponse{Documents: []*ai.Document{}}, nil
			}

			// Generate embedding for the query
			embeddingResp, err := genkit.Embed(ctx, g,
				ai.WithEmbedder(embedder),
				ai.WithTextDocs(queryText),
			)
			if err != nil {
				return nil, err
			}

			// Create a VectorRecord from the query embedding
			queryVector := VectorRecord{
				Prompt:    queryText,
				Embedding: embeddingResp.Embeddings[0].Embedding,
			}

			// Search for similar vectors using the MemoryVectorStore
			var similarRecords []VectorRecord
			var searchErr error

			if opts.MaxResults > 0 {
				// Use SearchTopNSimilarities for limited results
				similarRecords, searchErr = vectorStore.SearchTopNSimilarities(queryVector, opts.Limit, opts.MaxResults)
			} else {
				// Use SearchSimilarities for all results above threshold
				similarRecords, searchErr = vectorStore.SearchSimilarities(queryVector, opts.Limit)
			}

			if searchErr != nil {
				return nil, searchErr
			}

			// Convert VectorRecord results to ai.Document
			documents := make([]*ai.Document, len(similarRecords))
			for i, record := range similarRecords {
				// Create document with the prompt content and similarity score
				doc := ai.DocumentFromText(record.Prompt, map[string]any{
					"id":               record.Id,
					"cosine_similarity": record.CosineSimilarity,
				})
				documents[i] = doc
			}

			return &ai.RetrieverResponse{
				Documents: documents,
			}, nil
		},
	)
}