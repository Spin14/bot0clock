package model

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
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

type User struct {
	Username string `json:"username"`
}

func fromModel(model *userModel) (*User, error) {
	if model == nil {
		return nil, errors.New("nil Model")
	}
	if model.ID == 0 {
		return nil, errors.New("no ID Model")
	}

	return &User{Username:model.Username}, nil
}


func UsersPopulate(w http.ResponseWriter, r *http.Request) {
	db, closeDb := ProdStorage().dbCon()
	defer closeDb(db)

	var users []userModel
	db.Find(&users).Delete(userModel{})

	db.Create(&userModel{Username: "spin14"})
	db.Create(&userModel{Username: "chi"})

	users = []userModel{}
	db.Find(&users)

	_, _ = fmt.Fprint(w, "DB populated !\n")
}


type Storage struct {
	c connector
}

func (s Storage) dbCon() (*gorm.DB, func(db *gorm.DB)) {
	return s.c.GetCon()
}


func (s *Storage) Count() int {
	db, tearDown := s.dbCon()
	defer tearDown(db)

	var userModels []userModel
	db.Find(&userModels)
	return len(userModels)
}

func (s *Storage) Create(username string) (*User, error) {
	db, tearDown := s.dbCon()
	defer tearDown(db)

	instance := userModel{Username:username}
	db.Create(&instance)
	return fromModel(&instance)
}

func (s *Storage) ListAll() ([]User, error) {
	db, tearDown := s.dbCon()
	defer tearDown(db)

	var userModels []userModel
	db.Find(&userModels)

	var instances = make([]User, len(userModels))
	for i := range userModels {
		instance, err := fromModel(&userModels[i])
		if err != nil {
			return []User{}, err
		}
		instances[i] = *instance
	}

	return instances, nil
}


func (s *Storage) Retrieve(username string) (*User, error) {
	db, tearDown := s.dbCon()
	defer tearDown(db)

	var instance = userModel{}
	db.Last(&instance, "Username = ?", username)
	return fromModel(&instance)
}

func (s *Storage) Update(username string, newUsername string) (*User, error) {
	db, tearDown := s.dbCon()
	defer tearDown(db)

	var instance = userModel{}
	db.Last(&instance, "Username = ?", username)
	if instance.ID == 0 {
		return nil, errors.New("no ID Model")
	}

	db.Model(&instance).Update("Username", newUsername)
	return fromModel(&instance)
}
