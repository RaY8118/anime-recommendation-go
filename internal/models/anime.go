package models

type Title struct {
	Romaji         string `bson:"romaji,omitempty" json:"romaji,omitempty"`
	English        string `bson:"english,omitempty" json:"english,omitempty"`
	DisplayRomaji  string `bson:"display_romaji,omitempty" json:"display_romaji,omitempty"`
	DisplayEnglish string `bson:"display_english,omitempty" json:"display_english,omitempty"`
}

type CoverImage struct {
	Large string `bson:"large,omitempty" json:"large,omitempty"`
}

type Anime struct {
	ID           int        `bson:"id,omitempty" json:"id,omitempty"`
	Title        Title      `bson:"title,omitempty" json:"title"`
	Description  string     `bson:"description,omitempty" json:"description,omitempty"`
	Genres       []string   `bson:"genres,omitempty" json:"genres,omitempty"`
	AverageScore int        `bson:"averageScore,omitempty" json:"averageScore,omitempty"`
	Episodes     int        `bson:"episodes,omitempty" json:"episodes,omitempty"`
	Duration     int        `bson:"duration,omitempty" json:"duration,omitempty"`
	Season       string     `bson:"season,omitempty" json:"season,omitempty"`
	SeasonYear   int        `bson:"seasonYear,omitempty" json:"seasonYear,omitempty"`
	Status       string     `bson:"status,omitempty" json:"status,omitempty"`
	Source       string     `bson:"source,omitempty" json:"source,omitempty"`
	Studios      []string   `bson:"studios,omitempty" json:"studios,omitempty"`
	CoverImage   CoverImage `bson:"coverImage,omitempty" json:"coverImage,omitempty"`
	Embedding    []float32  `bson:"embedding,omitempty" json:"embedding,omitempty"`
}

type AnimeResponse struct {
	ID           int        `json:"id"`
	Title        Title      `json:"title"`
	Description  string     `json:"description"`
	Genres       []string   `json:"genres"`
	AverageScore int        `json:"averageScore"`
	Episodes     int        `json:"episodes"`
	Duration     int        `json:"duration"`
	Season       string     `json:"season"`
	SeasonYear   int        `json:"seasonYear"`
	Status       string     `json:"status"`
	Source       string     `json:"source"`
	Studios      []string   `json:"studios"`
	CoverImage   CoverImage `json:"coverImage"`
}

type AnimeTitleResponse struct {
	Romaji  string `json:"romaji"`
	English string `json:"english"`
}

type AnimeReccResponse struct {
	ID           int        `bson:"id,omitempty" json:"id"`
	Title        Title      `bson:"title,omitempty" json:"title"`
	Description  string     `bson:"description,omitempty" json:"description"`
	Genres       []string   `bson:"genres,omitempty" json:"genres"`
	AverageScore int        `bson:"averageScore,omitempty" json:"averageScore"`
	Episodes     int        `bson:"episodes,omitempty" json:"episodes"`
	Duration     int        `bson:"duration,omitempty" json:"duration"`
	Season       string     `bson:"season,omitempty" json:"season"`
	SeasonYear   int        `bson:"seasonYear,omitempty" json:"seasonYear"`
	Status       string     `bson:"status,omitempty" json:"status"`
	Source       string     `bson:"source,omitempty" json:"source"`
	Studios      []string   `bson:"studios,omitempty" json:"studios"`
	CoverImage   CoverImage `bson:"coverImage,omitempty" json:"coverImage"`
	Score        float64    `bson:"score,omitempty" json:"score"`
}

type AnimeAPIResponse struct {
	Page struct {
		Media []struct {
			ID    int `json:"id"`
			Title struct {
				Romaji  string `json:"romaji"`
				English string `json:"english"`
			} `json:"title"`
			Description  string   `json:"description"`
			Genres       []string `json:"genres"`
			AverageScore int      `json:"averageScore"`
			Episodes     int      `json:"episodes"`
			Duration     int      `json:"duration"`
			Season       string   `json:"season"`
			SeasonYear   int      `json:"seasonYear"`
			Status       string   `json:"status"`
			Source       string   `json:"source"`
			Studios      struct {
				Nodes []struct {
					Name string `json:"name"`
				} `json:"nodes"`
			} `json:"studios"`
			CoverImage struct {
				Large string `json:"large"`
			} `json:"coverImage"`
		} `json:"media"`
	} `json:"Page"`
}

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type OllamaResponse struct {
	Embeddings []float32 `json:"embedding"`
}
