package main

import (
	"anime/internal/database"
	"anime/internal/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	database.InitMongoDB()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	v1r := chi.NewRouter()

	v1r.Route("/anime", func(animeRouter chi.Router) {
		animeRouter.Get("/recommend", handlers.RecommendHandler)
		animeRouter.Get("/new-recommend", handlers.NewRecommendHandler)
		animeRouter.Get("/", handlers.AnimeByNameHandler)
		animeRouter.Get("/list", handlers.AnimeListHandler)
		animeRouter.Get("/random", handlers.RandomAnimeHandler)
		animeRouter.Get("/top-rated", handlers.TopRatedAnimesHandler)
		animeRouter.Get("/graphql", handlers.GraphQLAPIHandler)
		animeRouter.Get("/insert", handlers.InsertAnimeHandler)
		animeRouter.Get("/insertconcurrent", handlers.InsertAnimeConcurrentHandler)
	})

	r.Mount("/v1", v1r)

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}
