package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/spin14/bot0clock/model"
	"net/http"
)


type HttpHandler func(w http.ResponseWriter, r *http.Request)

func RetrieveUser(s *model.Storage) HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		username := mux.Vars(r)["username"]

		user, err := s.Retrieve(username)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(user)
	}
}

func CreateUser(s *model.Storage) HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		type createData struct {
			Username string `json:"username"`
		}

		var data = createData{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&data); err != nil {
			http.Error(w, "validation error", http.StatusBadRequest)
			return
		}


		user, err := s.Create(data.Username)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(user)
	}
}

func ListUsers(s *model.Storage) HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.ListAll()
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(users)
	}
}

func UpdateUser(s *model.Storage) HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		type updateData struct {
			Username string `json:"username"`
		}

		var data = updateData{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&data); err != nil {
			http.Error(w, "validation error", http.StatusBadRequest)
			return
		}

		username := mux.Vars(r)["username"]
		user, err := s.Update(username, data.Username)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(user)
	}
}
