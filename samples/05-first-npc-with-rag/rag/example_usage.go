package rag

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// ExampleUsage demonstrates how to use the custom retriever
func ExampleUsage(g *genkit.Genkit, embedder ai.Embedder) {
	ctx := context.Background()
	// Initialize the MemoryVectorStore
	vectorStore := &MemoryVectorStore{
		Records: make(map[string]VectorRecord),
	}

	// Create the memory vector retriever
	memoryRetriever := DefineMemoryVectorRetriever(g, vectorStore, embedder)

	// Example: Add some documents to the vector store first
	// AddDocumentToVectorStore(ctx, g, vectorStore, embedder, "Machine learning is a subset of AI")
	// AddDocumentToVectorStore(ctx, g, vectorStore, embedder, "Deep learning uses neural networks")

	// Create a query document
	queryDoc := ai.DocumentFromText("What is machine learning?", nil)

	// Create a retriever request with custom options
	request := &ai.RetrieverRequest{
		Query: queryDoc,
		Options: MemoryVectorRetrieverOptions{
			Limit:      0.7, // Higher similarity threshold
			MaxResults: 3,   // Return top 3 results
		},
	}

	// Use the memory vector retriever directly
	response, err := memoryRetriever.Retrieve(ctx, request)
	if err != nil {
		fmt.Printf("Error retrieving documents: %v\n", err)
		return
	}

	// Process the results
	fmt.Printf("Found %d documents:\n", len(response.Documents))
	for i, doc := range response.Documents {
		similarity := doc.Metadata["cosine_similarity"]
		id := doc.Metadata["id"]
		fmt.Printf("%d. ID: %s, Similarity: %v\n", i+1, id, similarity)
		fmt.Printf("   Content: %s\n", doc.Content[0].Text)
	}
}

// AddDocumentToVectorStore is a helper function to add documents to the vector store
func AddDocumentToVectorStore(ctx context.Context, g *genkit.Genkit, vectorStore *MemoryVectorStore, embedder ai.Embedder, text string) error {
	// Generate embedding for the text
	embeddingResp, err := genkit.Embed(ctx, g,
		ai.WithEmbedder(embedder),
		ai.WithTextDocs(text),
	)
	if err != nil {
		return err
	}

	// Create a VectorRecord
	record := VectorRecord{
		Prompt:    text,
		Embedding: embeddingResp.Embeddings[0].Embedding,
	}

	// Save to vector store
	_, err = vectorStore.Save(record)
	return err
}