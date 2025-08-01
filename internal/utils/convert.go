package utils

import "anime/internal/models"

func ConvertResponseToAnime(resp models.AnimeResponse, embedding []float32) models.Anime {
	return models.Anime{
		ID:           resp.ID,
		Title:        resp.Title,
		Description:  resp.Description,
		Genres:       resp.Genres,
		AverageScore: resp.AverageScore,
		Episodes:     resp.Episodes,
		Duration:     resp.Duration,
		Season:       resp.Season,
		SeasonYear:   resp.SeasonYear,
		Status:       resp.Status,
		Source:       resp.Source,
		Studios:      resp.Studios,
		CoverImage:   resp.CoverImage,
		Embedding:    embedding,
	}
}
