package main

type user struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	HashPassword string `json:"hash_password"`
}

func newUser(email, hashPassword string) *user {
	return &user{Email: email, HashPassword: hashPassword}
}
