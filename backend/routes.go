package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

func startServer(db *storage) {
	router := mux.NewRouter()
	router.HandleFunc("/accounts", createHandlerFunc(handleAccount, db))
	router.HandleFunc("/accounts/{email}", validateWithJWT(createHandlerFunc(handleAccountByEmail, db)))
	router.HandleFunc("/login", createHandlerFunc(handleLogin, db))

	log.Fatal(http.ListenAndServe(db.port, router))
}

func handleLogin(writer http.ResponseWriter, request *http.Request, db *storage) error {
	if request.Method != "POST" {
		return fmt.Errorf("invalid method")
	}

	userRequest := new(userRequest)
	if err := json.NewDecoder(request.Body).Decode(userRequest); err != nil {
		return err
	}

	userInDB, err := db.getUserByEmail(userRequest.Email)
	if err != nil {
		return err
	}

	if comparePassword(userRequest.Password, userInDB.HashPassword) {
		return fmt.Errorf("password %v incorrect", userRequest.Password)
	}

	//password is correct and we found the user in the db -> we can now give jwt key
	claims := &jwt.RegisteredClaims{ID: strconv.Itoa(userInDB.Id), ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24))}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(configuration.JWT_SECRET))

	if err != nil {
		return err
	}

	http.SetCookie(writer, &http.Cookie{Name: "jwt-token", Value: string(tokenString)})
	return WriteJSON(writer, http.StatusOK, "successfully logged in user and returned jwt as cookie")
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
	userRequest := new(userRequest)
	if err := json.NewDecoder(request.Body).Decode(userRequest); err != nil {
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
	claims := &jwt.RegisteredClaims{ID: strconv.Itoa(userObj.Id), ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24))}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(configuration.JWT_SECRET))

	if err != nil {
		return err
	}

	http.SetCookie(writer, &http.Cookie{Name: "jwt-token", Value: string(tokenString)})

	return WriteJSON(writer, http.StatusOK, userObj)
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

// validateWithJWT requires the request body to contain an "id" field for the user id we are trying to validate
func validateWithJWT(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("Validating with JWT")

		var tokenString string = ""
		for _, cookie := range request.Cookies() {
			if cookie.Name == "jwt-token" {
				tokenString = cookie.Value
			}
		}
		if tokenString == "" {
			WriteJSON(writer, http.StatusBadRequest, "cookie jwt-token does not exist")
			return
		}

		claims := &jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims,
			func(t *jwt.Token) (interface{}, error) {
				return []byte(configuration.JWT_SECRET), nil
			})

		if err != nil || !token.Valid {
			WriteJSON(writer, http.StatusForbidden, "access denied 1")
			return
		}

		idStruct := new(idStruct)
		if err := json.NewDecoder(request.Body).Decode(idStruct); err != nil {
			WriteJSON(writer, http.StatusForbidden, err)
			return
		}

		if claimID, err := strconv.Atoi(claims.ID); err != nil || idStruct.Id != claimID {
			WriteJSON(writer, http.StatusForbidden, "access denied 2")
			return
		}

		fmt.Printf("Authenticated User with ID %v\n", idStruct.Id)

		handler(writer, request)
	}
}

func extractVariable(request *http.Request, name string) (string, error) {
	email, ok := mux.Vars(request)[name]
	if !ok {
		return "", fmt.Errorf("could not find variable '%v'", name)
	}
	return email, nil
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
