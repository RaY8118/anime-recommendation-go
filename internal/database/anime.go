package database

import (
	"anime/internal/embeddings"
	"anime/internal/models"
	"anime/internal/utils"
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAnimeList() ([]models.Anime, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := AnimeCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("mongo find error: %w", err)
	}
	defer cursor.Close(ctx)

	var animes []models.Anime
	for cursor.Next(ctx) {
		var anime models.Anime
		if err := cursor.Decode(&anime); err != nil {
			log.Println("decode error:", err)
			continue
		}
		animes = append(animes, anime)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return animes, nil
}

func GetAnimeByName(name string) (models.AnimeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	lowerName := strings.ToLower(name)
	filter := bson.M{
		"$or": []bson.M{
			{"title.romaji": lowerName},
			{"title.english": lowerName},
		},
	}

	var anime models.AnimeResponse
	err := AnimeCollection.FindOne(ctx, filter).Decode(&anime)
	if err != nil {
		log.Println("mongo findOne error:", err)
		return models.AnimeResponse{}, err
	}
	return anime, nil
}

func GetRandomAnime() (models.AnimeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var anime models.AnimeResponse
	count, err := AnimeCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return models.AnimeResponse{}, err
	}

	if count == 0 {
		return models.AnimeResponse{}, fmt.Errorf("no animes found")
	}

	randomIndex := rand.Intn(int(count))

	cursor, err := AnimeCollection.Find(ctx, bson.M{}, options.Find().SetSkip(int64(randomIndex)).SetLimit(1))
	if err != nil {
		return models.AnimeResponse{}, fmt.Errorf("mongo find error: %w", err)
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if err := cursor.Decode(&anime); err != nil {
			return models.AnimeResponse{}, fmt.Errorf("decode error: %w", err)
		}
		return anime, nil
	}
	return models.AnimeResponse{}, fmt.Errorf("no animes found")

}

func GetTopRatedAnimes(limit int64) ([]models.AnimeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var animes []models.AnimeResponse
	cursor, err := AnimeCollection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"averageScore": -1}).SetLimit(limit))
	if err != nil {
		return animes, fmt.Errorf("mongo find error: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var anime models.AnimeResponse
		if err := cursor.Decode(&anime); err != nil {
			return animes, fmt.Errorf("decode error: %w", err)
		}
		animes = append(animes, anime)
	}

	return animes, nil

}
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
		_, err = NewAnimeCollection.InsertMany(ctx, animeDocs)
		if err != nil {
			return 0, fmt.Errorf("insert many error: %w", err)
		}
		log.Printf("Successfully inserted %d anime entries\n", len(animeDocs))
	} else {
		log.Println("No valid anime to insert")
	}
	return len(animeDocs), nil
}
