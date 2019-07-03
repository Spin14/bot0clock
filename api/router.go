package api

import (
	"github.com/gorilla/mux"
	"github.com/spin14/bot0clock/model"
)

func Router(s *model.Storage) *mux.Router {
	const usernamePattern = "{username:[A-Za-z0-9_]+}"

	router := mux.NewRouter()
	router.HandleFunc("/", ListUsers(s)).Methods("GET")
	router.HandleFunc("/", CreateUser(s)).Methods("POST")
	router.HandleFunc("/"+usernamePattern, RetrieveUser(s)).Methods("GET")
	router.HandleFunc("/"+usernamePattern, UpdateUser(s)).Methods("PUT")

	return router
}
