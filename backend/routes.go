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
	router.HandleFunc("/", createHandlerFunc(handleAccount))

	log.Fatal(http.ListenAndServe(db.port, router))
}

func handleAccount(writer http.ResponseWriter, request *http.Request) error {
	switch request.Method {
	case "POST":
		return handleCreateAccount(writer, request)
	case "GET":
		break
	}

	return fmt.Errorf("invalid method")
}

func handleCreateAccount(writer http.ResponseWriter, request *http.Request) error {
	return WriteJSON(writer, http.StatusOK, "message success")
}

func createHandlerFunc(handler func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if err := handler(writer, request); err != nil {
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
