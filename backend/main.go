package main

import (
	"os"
)

func main() {
	database := &storage{port: os.Getenv("DATABASE_URL")} //make sure to export this export DATABASE_URL=":8080"
	database.connectDb()
	defer database.database.Close()
	startServer(database)
}
