package utils

import (
	"anime/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GraphQLAPIRequest(page int64, perPage int64) ([]models.AnimeResponse, error) {
	query := `
	query ($page: Int, $perPage: Int) {
	  Page(page: $page, perPage: $perPage) {
		media(type: ANIME, sort: POPULARITY_DESC) {
		  id
		  title {
			romaji
			english
		  }
		  description
		  genres
		  averageScore
		  episodes
		  duration
		  season
		  seasonYear
		  status
		  source
		  studios {
			nodes {
			  name
			}
		  }
		  coverImage {
			large
		  }
		}
	  }
	}`

	variables := map[string]any{
		"page":    page,
		"perPage": perPage,
	}

	reqBody := models.GraphQLRequest{
		Query:     query,
		Variables: variables,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("https://graphql.anilist.co", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("anilist api error: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data models.AnimeAPIResponse `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var animes []models.AnimeResponse
	for _, m := range result.Data.Page.Media {
		studios := make([]string, 0, len(m.Studios.Nodes))
		for _, s := range m.Studios.Nodes {
			studios = append(studios, s.Name)
		}

		animes = append(animes, models.AnimeResponse{
			ID:           m.ID,
			Title:        models.Title{Romaji: m.Title.Romaji, English: m.Title.English},
			Description:  m.Description,
			Genres:       m.Genres,
			AverageScore: m.AverageScore,
			Episodes:     m.Episodes,
			Duration:     m.Duration,
			Season:       m.Season,
			SeasonYear:   m.SeasonYear,
			Status:       m.Status,
			Source:       m.Source,
			Studios:      studios,
			CoverImage:   models.CoverImage{Large: m.CoverImage.Large},
		})
	}

	return animes, nil
}
