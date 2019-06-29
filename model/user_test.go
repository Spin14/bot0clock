package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"testing"
)


func dbTestSetUp() (*gorm.DB, func(db *gorm.DB)) {
	db, dbTearDownF := dbConBase("test.db")
	db.AutoMigrate(&userModel{})
	assertNoUsers(db)
	return db, dbTestTearDown(dbTearDownF)
}

func assertNoUsers(db *gorm.DB) {
	var users []userModel
	db.Find(&users)
	if len(users) != 0 {
		panic("test db has records")
	}
}

func dbTestTearDown(f func(db *gorm.DB)) func(*gorm.DB) {
	return func(db *gorm.DB) {
		var users []userModel
		db.Find(&users).Delete(userModel{})
		f(db)
	}
}

func isUserF(reference *userModel) func(*userModel) (bool, string) {
	return func(candidate *userModel) (b bool, s string) {
		return reference.ID == candidate.ID, fmt.Sprintf(
			"referece user: %d (%s); candidate user: %d (%s)",
			reference.ID, reference.Username, candidate.ID, candidate.Username)
	}
}

func TestOrm(t *testing.T) {
	db, tearDownF := dbTestSetUp()
	defer tearDownF(db)

	var users []userModel
	var user userModel

	var spin14 = userModel{Username: "spin14"}
	var chi = userModel{Username: "chi"}

	isSpin14 := isUserF(&spin14)
	isChi := isUserF(&chi)

	users = []userModel{}
	db.Find(&users)
	if len(users) != 0 {
		t.Error("user count error")
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

	users = []userModel{}
	db.Find(&users)
	if len(users) != 2 {
		t.Error("user count error")
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
