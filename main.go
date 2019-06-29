package main

import (
	"github.com/gorilla/mux"
	"github.com/spin14/bot0clock/model"
	"log"
	"net/http"
)


func main() {
	model.DbSetUp()

	r := mux.NewRouter()

	usernamePattern := "{username}[A-Za-z0-9_]+"
	// Routes consist of a path and a handler function.
	r.HandleFunc("/users-populate", model.UsersPopulate).Methods("POST")
	r.HandleFunc("/users", model.UserList).Methods("GET")
	r.HandleFunc("/users", model.UserCreate).Methods("POST")
	r.HandleFunc("/users/" + usernamePattern, model.UserGet).Methods("GET")
	r.HandleFunc("/users/" + usernamePattern, model.UserUpdate).Methods("PUT")
	// Bind to a port and pass our router in
	log.Println("Started http server :: port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
