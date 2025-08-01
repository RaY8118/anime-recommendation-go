package embeddings

import (
	"context"
	"log"

	"google.golang.org/genai"
)

func GenerateEmbedding(query string) ([]float32, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	contents := []*genai.Content{
		genai.NewContentFromText(query, genai.RoleUser),
	}

	result, err := client.Models.EmbedContent(ctx,
		"gemini-embedding-001",
		contents,
		nil)

	if err != nil {
		log.Fatal(err)
	}

	return result.Embeddings[0].Values, nil
}
