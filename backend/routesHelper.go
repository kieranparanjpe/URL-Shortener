package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

// createJwtToken create a jwt with the id given. It then writes the token to a cookie via the writer provided.
func createJwtToken(writer http.ResponseWriter, id int) error {
	claims := &jwt.RegisteredClaims{ID: strconv.Itoa(id), ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24))}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(configuration.JWT_SECRET))

	if err != nil {
		return err
	}

	http.SetCookie(writer, &http.Cookie{Name: "jwt-token", Value: string(tokenString)})
	return nil
}

// createHandlerFunc should wrap all other handler funcs so that they can handle errors.
func createHandlerFunc(handler func(http.ResponseWriter, *http.Request, *storage) error, db *storage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if err := handler(writer, request, db); err != nil {
			writeJSON(writer, http.StatusBadRequest, jsonError{Error: err.Error()})
		}
	}
}

// getCookie gets the specified cookie from the request. if the cookie does not exist, returns a non-nil error
func getCookie(request *http.Request, name string) (string, error) {
	for _, cookie := range request.Cookies() {
		if cookie.Name == name {
			return cookie.Value, nil
		}
	}
	return "", fmt.Errorf("cookie %v does not exist", name)
}

// validateWithJWT requires the request body to contain an "id" field for the user id we are trying to validate
func validateWithJWT(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("\n----Running Middleware----")
		fmt.Println("Attempting to authenticate via jwt")

		tokenString, err := getCookie(request, "jwt-token")
		if err != nil {
			writeJSON(writer, http.StatusBadRequest, jsonError{Error: err.Error()})
			return
		}

		claims := &jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims,
			func(t *jwt.Token) (interface{}, error) {
				return []byte(configuration.JWT_SECRET), nil
			})

		if err != nil || !token.Valid {
			writeJSON(writer, http.StatusForbidden, jsonError{Error: "access denied"})
			return
		}

		idStruct := new(idStruct)
		if err := parseBody(request, idStruct, true); err != nil {
			writeJSON(writer, http.StatusForbidden, jsonError{Error: "access denied"})
			return
		}

		if claimID, err := strconv.Atoi(claims.ID); err != nil || idStruct.Id != claimID {
			writeJSON(writer, http.StatusForbidden, jsonError{Error: "access denied"})
			return
		}

		fmt.Printf("Authenticated User with ID %v\n", idStruct.Id)
		fmt.Println("----Ending Middleware----")

		handler(writer, request)
	}
}

func validateWithAdmin(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("\n----Running Middleware----")
		fmt.Println("Validating as admin")

		adminStruct := new(adminStruct)

		if err := parseBody(request, adminStruct, true); err != nil {
			writeJSON(writer, http.StatusForbidden, jsonError{Error: err.Error()})
			return
		}

		if adminStruct.AdminPassword != configuration.ADMIN_PASSWORD {
			writeJSON(writer, http.StatusForbidden, jsonError{Error: "admin password incorrect"})
			return
		}
		fmt.Println("Validated as admin")
		fmt.Println("----Ending Middleware----")

		handler(writer, request)
	}
}

/*
parseBody parses the json body in request into the interfact provided by out.
I think 'out' should be an explicit pointer but go doesn't like *interface{} and it works rn
resetReader can be set to true to allow for this request to be read multiple times IE if we are using it in middleware.
*/
func parseBody(request *http.Request, out interface{}, resetReader bool) error {
	if resetReader {
		readBuf := bytes.Buffer{}
		readBuf.ReadFrom(request.Body)
		backBuf := bytes.Buffer{}

		request.Body = io.NopCloser(io.TeeReader(&readBuf, &backBuf))

		defer func() {
			//request.Body = io.NopCloser(io.TeeReader(request.Body, &b))

			io.Copy(&readBuf, &backBuf)
		}()

	}
	if err := json.NewDecoder(request.Body).Decode(out); err != nil {
		return err
	}
	return nil
}

// extractVariable will extract and return a string variable from a request param. returns an error on failure.
func extractVariable(request *http.Request, name string) (string, error) {
	email, ok := mux.Vars(request)[name]
	if !ok {
		return "", fmt.Errorf("could not find variable '%v'", name)
	}
	return email, nil
}

// writeJSON writes a json message to the specified writer
func writeJSON(writer http.ResponseWriter, status int, value any) error {
	writer.WriteHeader(status)
	writer.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(writer).Encode(value)
}

type jsonError struct {
	Error string `json:"error"`
}

type jsonMessage struct {
	Message string `json:"message"`
}
