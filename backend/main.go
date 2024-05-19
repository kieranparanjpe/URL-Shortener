package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var configuration *config = &config{}

func main() {

	loadEnvFile()

	database := &storage{port: configuration.DATABASE_URL} //make sure to export this export DATABASE_URL=":8080"
	database.connectDb()
	defer database.database.Close()
	startServer(database)
}

func loadEnvFile() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
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
}

type config struct {
	DATABASE_URL string
	JWT_SECRET   string
}
