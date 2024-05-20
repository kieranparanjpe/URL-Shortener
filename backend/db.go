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
		password VARCHAR(500)
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

func (db *storage) dropUserById(id int) error {
	_, err := db.database.Query(`DELETE FROM "user" WHERE id=$1;`, id)
	return err
}

func (db *storage) dropAllLinks() error {
	_, err := db.database.Query(`DELETE TABLE IF EXISTS "link" CASCADE;`)
	return err
}

func (db *storage) addUser(u *user) error {
	rows, err := db.database.Query(`
	INSERT INTO "user"
	(email, password) VALUES
	($1, $2) 
	RETURNING id
	`, u.Email, u.HashPassword)

	if err != nil {
		return err
	}

	if rows.Next() {
		rows.Scan(&(u.Id))
	}

	return nil
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

func (db *storage) getUserById(id int) (*user, error) {
	rows, err := db.database.Query(`SELECT * FROM "user" WHERE id=$1`, id)
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

	return nil, fmt.Errorf("could not find account with id %v", id)
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

func (db *storage) addLink(l *link) error {
	rows, err := db.database.Query(`
	INSERT INTO "link"
	(url_redirect, user_id) VALUES
	($1, $2) 
	RETURNING id
	`, l.UrlRedirect, l.UserId)

	if err != nil {
		return err
	}

	if rows.Next() {
		rows.Scan(&(l.Id))
	}

	return nil
}

func (db *storage) getLinksByUserId(id int) ([](*link), error) {
	rows, err := db.database.Query(`
	SELECT link.id, url_redirect, user_id 
	FROM link JOIN "user" on user_id="user".id
	where user_id=$1
	`, id)

	if err != nil {
		return nil, err
	}

	links := make([](*link), 0)

	for rows.Next() {
		link := new(link)
		if err = rows.Scan(&link.Id, &link.UrlRedirect, &link.UserId); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}

func (db *storage) getLinkRedirect(linkID int) (urlRedirect string, err error) {
	rows, err := db.database.Query(`SELECT url_redirect FROM "link" WHERE id=$1`, linkID)
	if err != nil {
		return "", err
	}

	for rows.Next() {
		if err = rows.Scan(&urlRedirect); err != nil {
			return "", err
		}
		return urlRedirect, nil
	}
	return "", fmt.Errorf("could not find link with id %v", linkID)
}
