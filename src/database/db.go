package database

import (
	"ambassador/src/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {

	var err error

	DB, err = gorm.Open(mysql.Open("root:root@tcp(db:3306)/ambassador"), &gorm.Config{})

	if err != nil {
		panic("Could not connect database")
	}
}

func AutoMigrate() {
	DB.AutoMigrate(model.User{}, model.Product{}, model.Link{}, model.Order{}, model.OrderItem{})
}
