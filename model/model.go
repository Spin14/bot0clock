package model

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"math/rand"
	"net/http"
	"time"
)

type userModel struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username string  `gorm:"unique"`
	Salt     *string `gorm:"not null"`
	Password string
}


const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const saltRandLen = 10

func generateSalt(u *userModel) string {
	b := make([]byte, saltRandLen)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return fmt.Sprintf("%s%s", u.Username, string(b))
}

func (u *userModel) BeforeSave() (err error) {
	salt := generateSalt(u)
	u.Salt = &salt
	return nil
}

func UsersPopulate(w http.ResponseWriter, r *http.Request) {
	db, close := ProdStorage().dbCon()
	defer close(db)

	var users []userModel
	db.Find(&users).Delete(userModel{})

	db.Create(&userModel{Username: "spin14"})
	db.Create(&userModel{Username: "chi"})

	users = []userModel{}
	db.Find(&users)

	_, _ = fmt.Fprint(w, "DB populated !\n")
}

func UserCreate(w http.ResponseWriter, r *http.Request) {
	db, close := ProdStorage().dbCon()
	defer close(db)

	var user userModel
	decoder := json.NewDecoder(r.Body)
	_ = decoder.Decode(&user)
	db.Create(&user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

func UserList(w http.ResponseWriter, r *http.Request) {
	db, close := ProdStorage().dbCon()
	defer close(db)

	var users []userModel
	db.Find(&users)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(users)
}

func UserGet(w http.ResponseWriter, r *http.Request) {
	db, close := ProdStorage().dbCon()
	defer close(db)

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
	db, close := ProdStorage().dbCon()
	defer close(db)

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
