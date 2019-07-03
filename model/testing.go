package model

import (
	"github.com/jinzhu/gorm"
)

type testConnector struct{}

func (c testConnector) GetCon() (*gorm.DB, func(db *gorm.DB)) {
	db, tearDown := dbConBase("test.db")
	db.LogMode(false)
	return db, tearDown
}

func testStorage() *Storage {
	return &Storage{c: testConnector{}}
}

func InitUserModelTable() *Storage {
	s := testStorage()
	db, tearDown := s.dbCon()
	defer tearDown(db)
	db.AutoMigrate(&userModel{})
	assertNoUsers(db)
	return s
}

func CleanUserModelTable() {
	s := testStorage()
	db, tearDown := s.dbCon()
	defer tearDown(db)

	var users []userModel
	db.Find(&users).Delete(userModel{})
	assertNoUsers(db)
}

func assertNoUsers(db *gorm.DB) {
	var users []userModel
	db.Find(&users)
	if count := len(users); count != 0 {
		panic("there are user records in the test db")
	}
}
