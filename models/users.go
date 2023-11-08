package models

import "gorm.io/gorm"

type Users struct {
	ID       uint    `gorm:"primary key;autoIncrement" json:"id"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
}

func MigrateUsers(db *gorm.DB) error {
	err := db.AutoMigrate(&Users{})
	return err
}
