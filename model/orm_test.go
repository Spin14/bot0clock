package model

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"testing"
)

func Test_Orm(t *testing.T) {
	InitUserModelTable()
	defer CleanUserModelTable()

	s := testStorage()
	db, tearDownF := s.dbCon()
	defer tearDownF(db)
	var user userModel

	var spin14 = userModel{Username: "spin14"}
	var chi = userModel{Username: "chi"}

	var userCountEquals = func(db *gorm.DB, expected int) error {
		var users []userModel
		db.Find(&users)
		if count := len(users); count != expected {
			return errors.New(fmt.Sprintf("user count: %d; expected: %d", count, expected))
		}
		return nil
	}

	var isUser = func(reference *userModel) func(*userModel) (bool, string) {
		return func(candidate *userModel) (b bool, s string) {
			return reference.ID == candidate.ID, fmt.Sprintf(
				"referece user: %d (%s); candidate user: %d (%s)",
				reference.ID, reference.Username, candidate.ID, candidate.Username)
		}
	}

	isSpin14 := isUser(&spin14)
	isChi := isUser(&chi)

	if err := userCountEquals(db, 0); err != nil {
		t.Error(err.Error())
		return
	}

	user = userModel{}
	db.First(&user)
	if user.ID != 0 {
		t.Error("user first error")
		return
	}

	db.Create(&spin14)
	db.Create(&chi)

	if err := userCountEquals(db, 2); err != nil {
		t.Error(err.Error())
		return
	}

	user = userModel{}
	db.First(&user)
	if ok, msg := isSpin14(&user); !ok {
		t.Error(msg)
		return
	}

	user = userModel{}
	db.Last(&user)
	if ok, msg := isChi(&user); !ok {
		t.Error(msg)
		return
	}

	user = userModel{}
	db.Last(&user, "Username = ?", spin14.Username)
	if ok, msg := isSpin14(&user); !ok {
		t.Error(msg)
		return
	}

	db.Model(&spin14).Update("Username", "acci14")
	if spin14.Username != "acci14" {
		t.Error("update error")
		return
	}

	db.Delete(&spin14)
	user = userModel{}
	db.First(&user)
	if ok, msg := isChi(&user); !ok {
		t.Error(msg)
		return
	}
}
