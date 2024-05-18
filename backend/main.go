package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database := &storage{port: os.Getenv("DATABASE_URL")} //make sure to export this export DATABASE_URL=":8080"
	database.connectDb()
	defer database.database.Close()
	startServer(database)
}
