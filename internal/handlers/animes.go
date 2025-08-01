package handlers

import (
	"anime/internal/database"
	"anime/internal/embeddings"
	"anime/internal/models"
	"anime/internal/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func AnimeByNameHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing name parameter", http.StatusBadRequest)
		return
	}

	anime, err := database.GetAnimeByName(name)
	if err != nil {
		http.Error(w, "Failed to fetch anime", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(anime)
}

func AnimeListHandler(w http.ResponseWriter, r *http.Request) {
	animes, err := database.GetAnimeList()
	if err != nil {
		http.Error(w, "Failed to fetch animes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animes)
}

func RandomAnimeHandler(w http.ResponseWriter, r *http.Request) {
	anime, err := database.GetRandomAnime()
	if err != nil {
		http.Error(w, "Failed to fetch anime", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(anime)
}

func TopRatedAnimesHandler(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		http.Error(w, "Missing limit parameter", http.StatusBadRequest)
		return
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		http.Error(w, "Failed to parse limit", http.StatusInternalServerError)
		return
	}

	animes, err := database.GetTopRatedAnimes(limit)
	if err != nil {
		http.Error(w, "Failed to fetch animes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animes)
}

func RecommendHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}

	queryEmbedding, err := embeddings.GenerateEmbedding(query)
	if err != nil {
		http.Error(w, "Failed to generate embedding", http.StatusInternalServerError)
		return
	}

	animes, err := database.GetAnimeList()
	if err != nil {
		http.Error(w, "Failed to fetch anime list", http.StatusInternalServerError)
		return
	}

	type ScoredAnime struct {
		Anime models.Anime
		Score float64
	}
	var results []ScoredAnime

	for _, anime := range animes {
		if len(anime.Embedding) != len(queryEmbedding) {
			continue
		}

		a := make([]float32, len(queryEmbedding))
		b := make([]float32, len(queryEmbedding))
		for i := range queryEmbedding {
			a[i] = float32(queryEmbedding[i])
			b[i] = float32(anime.Embedding[i])
		}

		similarity, err := utils.CosineSimilarity(a, b)
		if err != nil {
			log.Println("Cosine error:", err)
			continue
		}

		results = append(results, ScoredAnime{Anime: anime, Score: similarity})
	}

	if len(results) == 0 {
		log.Println("No results found")
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	top := min(len(results), 1)
	topAnimes := make([]models.AnimeResponse, top)

	for i := range top {
		a := results[i].Anime
		topAnimes[i] = models.AnimeResponse{
			ID:           a.ID,
			Title:        a.Title,
			Description:  a.Description,
			Genres:       a.Genres,
			AverageScore: a.AverageScore,
			Episodes:     a.Episodes,
			Duration:     a.Duration,
			Season:       a.Season,
			SeasonYear:   a.SeasonYear,
			Status:       a.Status,
			Source:       a.Source,
			Studios:      a.Studios,
			CoverImage:   a.CoverImage,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topAnimes)
}

func NewRecommendHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}

	queryEmbedding, err := embeddings.GenerateEmbedding(query)
	if err != nil {
		http.Error(w, "Failed to generate embedding", http.StatusInternalServerError)
		return
	}

	vectorSearchStage := bson.D{
		{Key: "$vectorSearch", Value: bson.D{
			{Key: "index", Value: "embeddings_vector_index"},
			{Key: "path", Value: "embedding"},
			{Key: "queryVector", Value: queryEmbedding},
			{Key: "numCandidates", Value: 100},
			{Key: "limit", Value: 2},
		}}}

	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "id", Value: 1},
			{Key: "title", Value: 1},
			{Key: "description", Value: 1},
			{Key: "genres", Value: 1},
			{Key: "averageScore", Value: 1},
			{Key: "episodes", Value: 1},
			{Key: "duration", Value: 1},
			{Key: "season", Value: 1},
			{Key: "seasonYear", Value: 1},
			{Key: "status", Value: 1},
			{Key: "source", Value: 1},
			{Key: "studios", Value: 1},
			{Key: "coverImage", Value: 1},
			{Key: "score", Value: bson.D{{Key: "$meta", Value: "vectorSearchScore"}}},
		}}}

	cursor, err := database.AnimeCollection.Aggregate(context.Background(), mongo.Pipeline{vectorSearchStage, projectStage})
	if err != nil {
		http.Error(w, "Failed to fetch anime", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var animes []models.AnimeReccResponse
	for cursor.Next(context.Background()) {
		var anime models.AnimeReccResponse
		if err := cursor.Decode(&anime); err != nil {
			log.Println("decode error:", err)
			continue
		}
		animes = append(animes, anime)
	}

	if len(animes) == 0 {
		log.Println("No results found")
	}
	sort.Slice(animes, func(i, j int) bool {
		return animes[i].Score > animes[j].Score
	})

	top := min(len(animes), 2)
	topAnimes := make([]models.AnimeResponse, top)

	for i := range top {
		a := animes[i]
		topAnimes[i] = models.AnimeResponse{
			ID:           a.ID,
			Title:        a.Title,
			Description:  a.Description,
			Genres:       a.Genres,
			AverageScore: a.AverageScore,
			Episodes:     a.Episodes,
			Duration:     a.Duration,
			Season:       a.Season,
			SeasonYear:   a.SeasonYear,
			Status:       a.Status,
			Source:       a.Source,
			Studios:      a.Studios,
			CoverImage:   a.CoverImage,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topAnimes)

}

func GraphQLAPIHandler(w http.ResponseWriter, r *http.Request) {
	perPageStr := r.URL.Query().Get("perPage")
	if perPageStr == "" {
		http.Error(w, "Missing perPage parameter", http.StatusBadRequest)
		return
	}

	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		http.Error(w, "Missing page parameter", http.StatusBadRequest)
		return
	}

	perPage, err := strconv.ParseInt(perPageStr, 10, 64)
	if err != nil {
		http.Error(w, "Failed to parse perPage", http.StatusInternalServerError)
		return
	}

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		http.Error(w, "Failed to parse page", http.StatusInternalServerError)
		return
	}

	animes, err := utils.GraphQLAPIRequest(page, perPage)
	if err != nil {
		http.Error(w, "Failed to fetch animes", http.StatusInternalServerError)
		return
	}

	animesTitle := make([]string, 0, len(animes))

	for _, anime := range animes {
		animeTitle := anime.Title.Romaji
		animesTitle = append(animesTitle, animeTitle)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(animesTitle)
}

func InsertAnimeHandler(w http.ResponseWriter, r *http.Request) {
	perPageStr := r.URL.Query().Get("perPage")
	if perPageStr == "" {
		http.Error(w, "Missing perPage parameter", http.StatusBadRequest)
		return
	}

	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		http.Error(w, "Missing page parameter", http.StatusBadRequest)
		return
	}

	perPage, err := strconv.ParseInt(perPageStr, 10, 64)
	if err != nil {
		http.Error(w, "Failed to parse perPage", http.StatusInternalServerError)
		return
	}

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		http.Error(w, "Failed to parse perPage", http.StatusInternalServerError)
		return
	}

	count, err := database.InsertAnimes(page, perPage)
	if err != nil {
		http.Error(w, "Failed to insert animes", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{"status": "success", "message": "Animes inserted successfully", "count": count})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
