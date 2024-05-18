package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type storage struct {
	port     string
	database *sql.DB
}

func (db *storage) connectDb() {
	var err error
	db.database, err = sql.Open("postgres", db.port)
	if err != nil {
		log.Fatal(err)
		return
	}
}
