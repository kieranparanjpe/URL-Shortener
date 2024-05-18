package main

import (
	"os"
)

func main() {
	database := &storage{port: os.Getenv("DATABASE_URL")}
	database.connectDb()
	defer database.database.Close()
	startServer(database)
}
