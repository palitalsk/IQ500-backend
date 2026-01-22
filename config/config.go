package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PineconeApiKey   string
	PineconeIndexURL string
	EmbedServiceURL  string
	MongoDBURI       string
	DatabaseName     string
	GeminiApiKey     string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or failed to load")
	}

	config := &Config{
		PineconeApiKey:   os.Getenv("PINECONE_API_KEY"),
		PineconeIndexURL: os.Getenv("PINECONE_INDEX_URL"),
		EmbedServiceURL:  os.Getenv("EMBED_SERVICE_URL"),
		MongoDBURI:       os.Getenv("MONGODB_URI"),
		DatabaseName:     os.Getenv("DB_NAME"),
		GeminiApiKey:     os.Getenv("GEMINI_API_KEY"),
	}

	return config
}
