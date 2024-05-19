package main

import (
	"golang.org/x/crypto/bcrypt"
)

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type user struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	HashPassword string `json:"hash_password"`
}

func comparePassword(plain, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}

func newUser(email, password string) (*user, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &user{Email: email, HashPassword: string(hashPassword)}, nil
}

type linkRequest struct {
	UrlRedirect string `json:"url_redirect"`
	UserId      int    `json:"user_id"`
}

type link struct {
	Id          int    `json:"id"`
	UrlRedirect string `json:"url_redirect"`
	UserId      int    `json:"user_id"`
}

func newLink(urlRedirect string, userId int) *link {
	return &link{UrlRedirect: urlRedirect, UserId: userId}
}
