package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	apiPort := os.Getenv("API_PORT")

	if rabbitMQURL == "" || apiPort == "" {
		log.Fatal("Missing required environment variables")
	}

}
