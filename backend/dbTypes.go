package main

import (
	"golang.org/x/crypto/bcrypt"
)

type idStruct struct {
	Id     int `json:"id"`
	LinkId int `json:"link_id"`
}

type adminStruct struct {
	AdminPassword string `json:"admin_password"`
}

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
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}

//nil == true
//!nil == false

func newUser(email, password string) (*user, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &user{Email: email, HashPassword: string(hashPassword)}, nil
}

func (u *user) noPassword() user {
	return user{Id: u.Id, Email: u.Email}
}

type linkRequest struct {
	UrlRedirect string `json:"url_redirect"`
	UserId      int    `json:"id"`
}

type link struct {
	Id          int    `json:"id"`
	UrlRedirect string `json:"url_redirect"`
	UserId      int    `json:"user_id"`
	EncodedId   string `json:"encoded_id"`
}

func newLink(urlRedirect string, userId int) *link {
	return &link{UrlRedirect: urlRedirect, UserId: userId}
}
