package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var configuration *config = &config{}

const urlKey int = 12951

func main() {
	loadEnvFile()

	database := &storage{port: configuration.DATABASE_URL}
	database.connectDb()
	defer database.database.Close()
	startServer(database)

}

func loadEnvFile() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Println("Error loading .env file")
	}

	var ok bool
	configuration.DATABASE_URL, ok = os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatal("could not find DATABASE_URL in environment variables")
	}

	configuration.JWT_SECRET, ok = os.LookupEnv("JWT_SECRET")
	if !ok {
		log.Fatal("could not find JWT_SECRET in environment variables")
	}

	configuration.ADMIN_PASSWORD, ok = os.LookupEnv("ADMIN_PASSWORD")
	if !ok {
		log.Fatal("could not find ADMIN_PASSWORD in environment variables")
	}

	configuration.POSTGRES_URL, ok = os.LookupEnv("POSTGRES_URL")
	if !ok {
		log.Fatal("could not find POSTGRES_URL in environment variables")
	}
}

type config struct {
	DATABASE_URL   string
	JWT_SECRET     string
	ADMIN_PASSWORD string
	POSTGRES_URL   string
}
