package model

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"net/http"
)

type userModel struct {
	gorm.Model
	Username string `json:"username"`
}

func DbSetUp() {
	db, f := dbCon()
	defer f(db)
	db.AutoMigrate(&userModel{})
}

func dbConBase(dbName string) (*gorm.DB, func(db *gorm.DB)) {
	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		panic("failed to connect database")
	}

	return db, func(db *gorm.DB) {
		if err := db.Close(); err != nil {
			panic("failed to close database")
		}
	}
}

func dbCon() (*gorm.DB, func(db *gorm.DB)) {
	return dbConBase("dev.db")
}

func UsersPopulate(w http.ResponseWriter, r *http.Request) {
	db, closeF := dbCon()
	defer closeF(db)

	var users []userModel
	db.Find(&users).Delete(userModel{})

	db.Create(&userModel{Username: "spin14"})
	db.Create(&userModel{Username: "chi"})

	users = []userModel{}
	db.Find(&users)

	_, _ = fmt.Fprint(w, "DB populated !\n",)
}

func UserCreate(w http.ResponseWriter, r *http.Request) {
	db, closeF := dbCon()
	defer closeF(db)

	var user userModel
	decoder := json.NewDecoder(r.Body)
	_ = decoder.Decode(&user)
	db.Create(&user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

func UserList(w http.ResponseWriter, r *http.Request) {
	db, closeF := dbCon()
	defer closeF(db)
	
	var users []userModel
	db.Find(&users)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(users)
}

func UserGet(w http.ResponseWriter, r *http.Request) {
	db, closeF := dbCon()
	defer closeF(db)

	username := mux.Vars(r)["username"]

	var user userModel
	db.Where("username = ?", username).First(&user)

	if user.ID == 0 {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(userModel{Username: username})
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	db, closeF := dbCon()
	defer closeF(db)

	username := mux.Vars(r)["username"]

	var user userModel
	db.Where("username = ?", username).First(&user)

	if user.ID == 0 {
		http.NotFound(w, r)
		return
	}

	decoder := json.NewDecoder(r.Body)
	_ = decoder.Decode(&user)

	db.Model(&user).Update("Username", user.Username)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

