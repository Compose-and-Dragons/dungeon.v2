package agents

import (
	"context"
	"fmt"
	"npc-agent-with-tools/msg"
	"npc-agent-with-tools/rag"

	"github.com/firebase/genkit/go/ai"
)

func retrieveSimilarDocuments(ctx context.Context, query string, retriever ai.Retriever, similarityThreshold float64, similarityMaxResults int) (string, error) {
	// Create a query document from the user question
	queryDoc := ai.DocumentFromText(query, nil)

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

	msg.DisplaySimilarityMessages(
		"--------------------------------------------------",
		fmt.Sprintf("\n📘 Found %d similar documents:", len(retrieveResponse.Documents)),
	)

	for i, doc := range retrieveResponse.Documents {
		similarity := doc.Metadata["cosine_similarity"].(float64)
		id := doc.Metadata["id"].(string)
		content := doc.Content[0].Text

		msg.DisplaySimilarityMessages(
			fmt.Sprintf("%d. ID: %s, Similarity: %.4f\n", i+1, id, similarity),
			fmt.Sprintf("   Content: %s\n\n", content),
		)

		similarDocuments += content
	}

	msg.DisplaySimilarityMessages(
		"--------------------------------------------------",
		"",
	)

	return similarDocuments, nil
}
