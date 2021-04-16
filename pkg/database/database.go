package database

import (
	"errors"

	"github.com/saman2000hoseini/virtual-box_management_server/internal/app/virtual-box/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("./myDB.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

//create tables from models
func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&model.User{})
	return err
}

//make database and create tables for the first time
func FirstSetup() (*gorm.DB, error) {
	db, err := NewDB()
	if err != nil {
		return nil, errors.New("error on creating db")
	}
	err = migrate(db)
	if err != nil {
		return nil, errors.New("error on creating tables" + err.Error())
	}
	return db, nil
}
