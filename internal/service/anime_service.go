package service

import (
	"anime/internal/database"
	"anime/internal/embeddings"
	"anime/internal/models"
	"anime/internal/utils"
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

func InsertAnimes(page int64, perPage int64) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	animes, err := utils.GraphQLAPIRequest(page, perPage)
	if err != nil {
		return 0, err
	}

	var animeDocs []any

	for _, animeResp := range animes {
		embeddingText := animeResp.Title.Romaji + " " + animeResp.Title.English + " " + animeResp.Description + " " + strings.Join(animeResp.Genres, " ")

		embedding, err := embeddings.GenerateEmbedding(embeddingText)
		if err != nil {
			log.Println("embedding error:", err)
			continue
		}

		anime := utils.ConvertResponseToAnime(animeResp, embedding)

		animeDocs = append(animeDocs, anime)
	}

	if len(animeDocs) > 0 {
		_, err = database.NewAnimeCollection.InsertMany(ctx, animeDocs)
		if err != nil {
			return 0, fmt.Errorf("insert many error: %w", err)
		}
		log.Printf("Successfully inserted %d anime entries\n", len(animeDocs))
	} else {
		log.Println("No valid anime to insert")
	}
	return len(animeDocs), nil
}

func InsertAnimesConcurrent(startPage int64, endPage int64, perPage int64) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var (
		wg            sync.WaitGroup
		mu            sync.Mutex
		allDocs       []any
		totalInserted int
	)

	type result struct {
		docs []any
		err  error
	}

	resultsChan := make(chan result)

	for page := startPage; page <= endPage; page++ {
		wg.Add(1)

		go func(p int64) {
			defer wg.Done()

			animes, err := utils.GraphQLAPIRequest(p, perPage)
			if err != nil {
				log.Printf("error fetching page %d: %v", p, err)
				resultsChan <- result{nil, err}
				return
			}

			var innerWg sync.WaitGroup
			var docs []any
			var innerMu sync.Mutex

			for _, animeResp := range animes {
				innerWg.Add(1)

				go func(ar models.AnimeResponse) {
					defer innerWg.Done()

					embeddingText := ar.Title.Romaji + " " + ar.Title.English + " " + ar.Description + " " + strings.Join(ar.Genres, " ")

					embedding, err := embeddings.GenerateEmbeddingsOllama(embeddingText)
					if err != nil {
						log.Println("embedding error:", err)
						return
					}

					doc := utils.ConvertResponseToAnime(ar, embedding)

					innerMu.Lock()
					docs = append(docs, doc)
					innerMu.Unlock()
				}(animeResp)
			}

			innerWg.Wait()
			resultsChan <- result{docs, nil}
		}(page)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for res := range resultsChan {
		if res.err != nil {
			continue
		}

		mu.Lock()
		allDocs = append(allDocs, res.docs...)
		totalInserted += len(res.docs)
		mu.Unlock()
	}

	if len(allDocs) > 0 {
		_, err := database.NewAnimeCollection.InsertMany(ctx, allDocs)
		if err != nil {
			return 0, fmt.Errorf("insert many error: %w", err)
		}
		log.Printf("Successfully inserted %d anime entries\n", totalInserted)
	} else {
		log.Println("No valid anime to insert")
	}

	return totalInserted, nil

}
