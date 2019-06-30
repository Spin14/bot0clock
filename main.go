package main

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/spin14/bot0clock/model"
	"log"
	"net/http"
	"strings"
)


func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func getToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no token
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("'Authorization' header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}


func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getToken(r)
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		if token != "hello" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}


func main() {
	model.MigrateUsersTable()

	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	usernamePattern := "{username}[A-Za-z0-9_]+"

	sr := r.PathPrefix("/users").Subrouter()
	sr.HandleFunc("/", model.UserList).Methods("GET")
	sr.HandleFunc("/" + usernamePattern, model.UserGet).Methods("GET")
	sr.HandleFunc("/" + usernamePattern, model.UserUpdate).Methods("PUT")
	sr.Use(authMiddleware)

	sr.HandleFunc("/", model.UserCreate).Methods("POST")


	// Routes consist of a path and a handler function.
	r.HandleFunc("/users-populate", model.UsersPopulate).Methods("POST")

	// Bind to a port and pass our router in
	log.Println("Started http server :: port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
