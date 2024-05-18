package main

type user struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	HashPassword string `json:"hash_password"`
}

func newUser(email, hashPassword string) *user {
	return &user{Email: email, HashPassword: hashPassword}
}

type link struct {
	Id          int    `json:"id"`
	UrlRedirect string `json:"url_redirect"`
	UserId      int    `json:"user_id"`
}
