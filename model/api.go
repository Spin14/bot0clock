package model

import (
	"errors"
	"github.com/jinzhu/gorm"
)

type User struct {
	Username string
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


type connector interface {
	GetCon() (*gorm.DB, func(db *gorm.DB))
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

func (s *Storage) ListAll() ([]User, error) {
	db, tearDown := s.dbCon()
	defer tearDown(db)

	var userModels []userModel
	db.Find(&userModels)

	var users = make([]User, len(userModels))
	for i := range userModels {
		user, err := fromModel(&userModels[i])
		if err != nil {
			return []User{}, err
		}
		users[i] = *user
	}

	return users, nil
}

func (s *Storage) Create(username string) (*User, error) {
	db, tearDown := s.dbCon()
	defer tearDown(db)

	instance := userModel{Username:username}
	db.Create(&instance)
	return fromModel(&instance)
}

