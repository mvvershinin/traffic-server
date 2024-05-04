package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	SrvAddress string
}

func Configure() Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//localhost:50051
	c := Config{
		SrvAddress: os.Getenv("SERVER_ADDRESS_PORT"),
	}

	return c
}
