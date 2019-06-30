package model

import "github.com/jinzhu/gorm"


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

type prodConnector struct{}


func (c prodConnector) GetCon() (*gorm.DB, func(db *gorm.DB)) {
	db, tearDown := dbConBase("dev.db")
	return db, tearDown
}

func ProdStorage() *Storage {
	return &Storage{c: prodConnector{}}
}


func MigrateUsersTable() {
	db, tearDown := ProdStorage().dbCon()
	defer tearDown(db)
	db.AutoMigrate(&userModel{})
}