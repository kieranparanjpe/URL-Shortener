package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func startServer(db *storage) {
	router := mux.NewRouter()
	router.HandleFunc("/accounts", createHandlerFunc(handleAccount, db))
	router.HandleFunc("/accounts/id", validateWithJWT(createHandlerFunc(handleAccountById, db)))
	router.HandleFunc("/login", createHandlerFunc(handleLogin, db))
	router.HandleFunc("/logout", createHandlerFunc(handleLogout, db))

	router.HandleFunc("/links", validateWithJWT(createHandlerFunc(handleLink, db)))
	router.Handle("/l/{id}", createHandlerFunc(handleFollowLink, db))

	log.Fatal(http.ListenAndServe(db.port, router))
}

func handleLogin(writer http.ResponseWriter, request *http.Request, db *storage) error {
	if request.Method != "POST" {
		return writeJSON(writer, http.StatusMethodNotAllowed, jsonError{Error: "invalid method provided"})
	}

	userRequest := new(userRequest)
	if err := parseBody(request, userRequest, false); err != nil {
		return err
	}

	userInDB, err := db.getUserByEmail(userRequest.Email)
	if err != nil {
		return err
	}

	if !comparePassword(userRequest.Password, userInDB.HashPassword) {
		return writeJSON(writer, http.StatusForbidden, jsonError{Error: fmt.Sprintf("password %v incorrect", userRequest.Password)})
	}
	//password is correct and we found the user in the db -> we can now give jwt key

	if err := createJwtToken(writer, userInDB.Id); err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, userInDB.noPassword())
}

func handleLogout(writer http.ResponseWriter, request *http.Request, db *storage) error {
	if request.Method != "POST" {
		return writeJSON(writer, http.StatusMethodNotAllowed, jsonError{Error: "invalid method provided"})
	}

	_, err := getCookie(request, "jwt-token")
	if err != nil {
		return writeJSON(writer, http.StatusBadRequest, jsonError{Error: err.Error()})
	}

	http.SetCookie(writer, &http.Cookie{Name: "jwt-token", Value: "", Expires: time.Now().Add(-10)})
	return writeJSON(writer, http.StatusOK, jsonMessage{Message: "successfully logged out user"})
}

func handleAccount(writer http.ResponseWriter, request *http.Request, db *storage) error {
	switch request.Method {
	case "POST":
		return handleCreateAccount(writer, request, db)
	case "GET":
		return handleGetAllAccounts(writer, request, db)
	case "DELETE":
		return handleDropAllAccounts(writer, request, db)
	}

	return writeJSON(writer, http.StatusMethodNotAllowed, jsonError{Error: "invalid method provided"})
}

func handleAccountById(writer http.ResponseWriter, request *http.Request, db *storage) error {
	switch request.Method {
	case "GET":
		return handleGetAccountById(writer, request, db)
	case "DELETE":
		return handleDropAccountById(writer, request, db)
	}

	return writeJSON(writer, http.StatusMethodNotAllowed, jsonError{Error: "invalid method provided"})
}

func handleCreateAccount(writer http.ResponseWriter, request *http.Request, db *storage) error {
	userRequest := new(userRequest)
	if err := parseBody(request, userRequest, false); err != nil {
		return err
	}
	userObj, err := newUser(userRequest.Email, userRequest.Password)
	if err != nil {
		return err
	}

	if err := db.addUser(userObj); err != nil {
		return err
	}

	//sign in the user now:
	if err := createJwtToken(writer, userObj.Id); err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, userObj.noPassword())
}

func handleGetAllAccounts(writer http.ResponseWriter, request *http.Request, db *storage) error {
	users, err := db.getAllUsers()
	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, users)
}

func handleGetAccountById(writer http.ResponseWriter, request *http.Request, db *storage) error {
	idStruct := new(idStruct)
	if err := parseBody(request, idStruct, false); err != nil {
		return err
	}

	user, err := db.getUserById(idStruct.Id)
	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, user.noPassword())
}

func handleDropAllAccounts(writer http.ResponseWriter, request *http.Request, db *storage) error {
	err := db.dropAllUsers()
	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, jsonMessage{Message: "successfully dropped all accounts from database"})
}

func handleDropAccountById(writer http.ResponseWriter, request *http.Request, db *storage) error {
	idStruct := new(idStruct)
	if err := parseBody(request, idStruct, false); err != nil {
		return err
	}

	err := db.dropUserById(idStruct.Id)
	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, jsonMessage{Message: fmt.Sprintf("successfully dropped account id=%v from database", idStruct.Id)})
}

func handleDropAllLinks(writer http.ResponseWriter, request *http.Request, db *storage) error {
	err := db.dropAllLinks()
	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, jsonMessage{Message: "successfully dropped all links from database"})
}

func handleLink(writer http.ResponseWriter, request *http.Request, db *storage) error {
	switch request.Method {
	case "POST":
		return handleCreateLink(writer, request, db)
	case "GET":
		return handleGetLinks(writer, request, db)
	case "DELETE":
		return handleDropAllAccounts(writer, request, db)
	}

	return writeJSON(writer, http.StatusMethodNotAllowed, jsonError{Error: "invalid method provided"})
}

func handleCreateLink(writer http.ResponseWriter, request *http.Request, db *storage) error {
	linkRequest := new(linkRequest)
	if err := parseBody(request, linkRequest, false); err != nil {
		return err
	}

	link := newLink(linkRequest.UrlRedirect, linkRequest.UserId)
	if err := db.addLink(link); err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, link)
}

func handleGetLinks(writer http.ResponseWriter, request *http.Request, db *storage) error {
	idStruct := new(idStruct)
	if err := parseBody(request, idStruct, false); err != nil {
		return err
	}

	links, err := db.getLinksByUserId(idStruct.Id)

	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, links)
}

func handleFollowLink(writer http.ResponseWriter, request *http.Request, db *storage) error {
	if request.Method != "GET" {
		return writeJSON(writer, http.StatusMethodNotAllowed, jsonError{Error: "invalid method provided"})
	}

	idString, err := extractVariable(request, "id")
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		return err
	}

	url, err := db.getLinkRedirect(id)

	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, map[string]string{"redirect_url": url})
}
