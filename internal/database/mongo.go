package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var AnimeCollection *mongo.Collection
var NewAnimeCollection *mongo.Collection

func InitMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	MongoDBEnv := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDBEnv))
	if err != nil {
		return
	}

	MongoClient = client
	AnimeCollection = MongoClient.Database("anime_recommendation").Collection("animes")
	NewAnimeCollection = MongoClient.Database("anime_recommendation").Collection("new_animes")

	log.Println("Connected to MongoDB")

}

func CloseMongoDB() {
	if MongoClient != nil {
		if err := MongoClient.Disconnect(context.Background()); err != nil {
			log.Println("Error disconnecting from MongoDB:", err)
		} else {
			log.Println("Disconnected from MongoDB")
		}
	}
}
