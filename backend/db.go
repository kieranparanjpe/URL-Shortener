package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type storage struct {
	port     string
	database *sql.DB
}

func (db *storage) connectDb() {
	connStr := "user=postgres dbname=postgres password=secret_postrgresql sslmode=disable"
	var err error
	db.database, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := db.database.Ping(); err != nil {
		log.Fatal(err)
	}
	if err := db.createTables(); err != nil {
		log.Fatal(err)
	}
}

func (db *storage) createTables() error {
	_, err := db.database.Query(`CREATE TABLE IF NOT EXISTS "user"
	(
		id SERIAL PRIMARY KEY,
		email VARCHAR(100) UNIQUE,
		password VARCHAR(50)
	);`)

	if err != nil {
		return err
	}

	_, err = db.database.Query(`CREATE TABLE IF NOT EXISTS "link"
	(
		id SERIAL PRIMARY KEY,
		url_redirect VARCHAR(1000),
		user_id INT,
		FOREIGN KEY (user_id) REFERENCES "user" (id)
	);`)

	if err != nil {
		return err
	}

	return nil
}

func (db *storage) dropAllUsers() error {
	_, err := db.database.Query(`DROP TABLE IF EXISTS "user" CASCADE;`)
	return err
}

func (db *storage) dropUserByEmail(email string) error {
	_, err := db.database.Query(`DELETE FROM "user" WHERE email=$1;`, email)
	return err
}

func (db *storage) dropAllLinks() error {
	_, err := db.database.Query(`DELETE TABLE IF EXISTS "link" CASCADE;`)
	return err
}

func (db *storage) addUser(u *user) error {
	_, err := db.database.Query(`
	INSERT INTO "user"
	(email, password) VALUES
	($1, $2)
	`, u.Email, u.HashPassword)

	return err
}

func (db *storage) getAllUsers() ([](*user), error) {
	rows, err := db.database.Query(`SELECT * FROM "user"`)

	if err != nil {
		return nil, err
	}
	users := make([](*user), 0)
	for rows.Next() {
		user := new(user)
		if err = rows.Scan(&user.Id, &user.Email, &user.HashPassword); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (db *storage) getUserByEmail(email string) (*user, error) {
	rows, err := db.database.Query(`SELECT * FROM "user" WHERE email=$1`, email)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user := new(user)
		if err = rows.Scan(&user.Id, &user.Email, &user.HashPassword); err != nil {
			return nil, err
		}
		return user, nil
	}

	return nil, fmt.Errorf("could not find account with email %v", email)
}
