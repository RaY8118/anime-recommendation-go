# ✨ NekoRec: Anime Recommendation System (Go Edition)

This project is a rewrite of my main anime recommendation system. NekoRec is a backend application designed to provide personalized anime recommendations. It allows users to browse anime, view details, explore by genre, and receive recommendations based on content similarity. This version features a backend built with Go.

## 🚀 Features

- 📚 **Browse Anime:** Explore a wide collection of anime.
- 🎬 **Anime Details:** View detailed information about each anime, including synopsis, genres, and more.
- 🏷️ **Genre Exploration:** Filter and discover anime by various genres.
- ❤️ **Personalized Recommendations:** Get anime recommendations based on content similarity (e.g., description, genres).

## 🛠️ Technologies Used

### 🐹 Backend (Go)

- **Go:** The core programming language.
- **Chi Web Framework:** A lightweight, idiomatic, and composable router for building HTTP services.
- **Air:** For live reloading during development.
- **MongoDB:** NoSQL database used for storing anime data and embeddings.
- **Anilist API:** Data for anime is sourced from the Anilist API.
- **Gemini Embeddings:** Used for generating vector embeddings of anime descriptions for content-based recommendations.
- **Gemini API:** Utilized for the AI chatbot functionality.
- **Data Processing:** Custom scripts for fetching anime data, text cleaning, generating embeddings, and calculating similarity for recommendations.

## 🏁 Getting Started

Follow these instructions to set up and run the project locally.

### 📋 Prerequisites

- **Go 1.18+** (or the latest stable version)
- **Git**
- **MongoDB** instance running (local or remote)

### 🐹 Backend Setup

1.  Navigate to the project root directory:
    ```bash
    cd /anime-recommendation-go
    ```
2.  Install Go dependencies:
    ```bash
    go mod tidy
    ```
3.  Install Air (if you haven't already):
    ```bash
    go install github.com/cosmtrek/air@latest
    ```
4.  Run the Go application with Air for live reloading:
    ```bash
    air
    ```
    The backend API will be available at `http://localhost:8080` (or similar, check `cmd/anime/main.go` for the exact port).

## 📂 Project Structure

```
.
├── .gitignore
├── go.mod
├── go.sum
├── cmd/                          # Main application entry points
│   └── anime/
│       ├── .air.toml             # Configuration for Air (live-reloading for Go apps)
│       └── main.go               # Main Go application entry point
├── internal/                     # Internal packages and business logic
│   ├── database/                 # Database interactions (MongoDB)
│   │   ├── anime.go
│   │   └── mongo.go
│   ├── embeddings/               # Logic for generating and handling embeddings
│   │   └── embed.go
│   ├── handlers/                 # HTTP request handlers
│   │   └── animes.go
│   ├── models/                   # Data models/structs
│   │   └── anime.go
│   ├── repository/               # Data access layer
│   ├── service/                  # Business logic layer
│   └── utils/                    # Utility functions
│       ├── convert.go
│       ├── graphql.go
│       └── similarity.go
└── tmp/                          # Temporary files/data
```

## 🤝 Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.
