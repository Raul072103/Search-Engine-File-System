package main

import (
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// Env file setup
	err := godotenv.Load("./../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}
