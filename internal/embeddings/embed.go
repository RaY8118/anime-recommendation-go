package embeddings

import (
	"anime/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

func GenerateEmbeddingsOllama(query string) ([]float32, error) {
	ollamaURL := os.Getenv("OLLAMA_URL")
	modelName := os.Getenv("OLLAMA_MODEL")

	reqBody := models.OllamaRequest{
		Model:  modelName,
		Prompt: query,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Post(ollamaURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned non-200 status: %s", resp.Status)
	}

	var res models.OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to parse ollama respone: %w", err)
	}

	return res.Embeddings, nil
}
