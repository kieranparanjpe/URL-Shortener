package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func startServer(db *storage) {
	router := mux.NewRouter()
	router.HandleFunc("/accounts", createHandlerFunc(handleAccount, db))
	router.HandleFunc("/accounts/{email}", createHandlerFunc(handleAccountByEmail, db))

	log.Fatal(http.ListenAndServe(db.port, router))
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

	return fmt.Errorf("invalid method")
}

func handleAccountByEmail(writer http.ResponseWriter, request *http.Request, db *storage) error {
	switch request.Method {
	case "GET":
		return handleGetAccountByEmail(writer, request, db)
	case "DELETE":
		return handleDropAccountByEmail(writer, request, db)
	}

	return fmt.Errorf("invalid method")
}

func handleCreateAccount(writer http.ResponseWriter, request *http.Request, db *storage) error {
	u := new(user)
	if err := json.NewDecoder(request.Body).Decode(u); err != nil {
		return err
	}
	if err := db.addUser(u); err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, u)
}

func handleGetAllAccounts(writer http.ResponseWriter, request *http.Request, db *storage) error {
	users, err := db.getAllUsers()
	if err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, users)
}

func handleGetAccountByEmail(writer http.ResponseWriter, request *http.Request, db *storage) error {
	email, err := extractVariable(request, "email")
	if err != nil {
		return err
	}

	u, err := db.getUserByEmail(email)
	if err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, u)
}

func handleDropAllAccounts(writer http.ResponseWriter, request *http.Request, db *storage) error {
	err := db.dropAllUsers()
	if err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, jsonMessage{Message: "successfully dropped all accounts from database"})
}

func handleDropAccountByEmail(writer http.ResponseWriter, request *http.Request, db *storage) error {
	email, err := extractVariable(request, "email")
	if err != nil {
		return err
	}

	err = db.dropUserByEmail(email)
	if err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, jsonMessage{Message: "successfully dropped account from database"})
}

func extractVariable(request *http.Request, name string) (string, error) {
	email, ok := mux.Vars(request)[name]
	if !ok {
		return "", fmt.Errorf("could not find variable '%v'", name)
	}
	return email, nil
}

func handleDropAllLinks(writer http.ResponseWriter, request *http.Request, db *storage) error {
	err := db.dropAllLinks()
	if err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, jsonMessage{Message: "successfully dropped all links from database"})
}

func createHandlerFunc(handler func(http.ResponseWriter, *http.Request, *storage) error, db *storage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if err := handler(writer, request, db); err != nil {
			WriteJSON(writer, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func WriteJSON(writer http.ResponseWriter, status int, value any) error {
	writer.WriteHeader(status)
	writer.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(writer).Encode(value)
}

type ApiError struct {
	Error string `json:"error"`
}

type jsonMessage struct {
	Message string `json:"message"`
}
