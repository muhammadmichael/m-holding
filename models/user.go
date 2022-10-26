package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `form:"name" json: "name" validate:"required"`
	Username string `form:"username" json: "username" validate:"required"`
	Image    string `form:"image" json: "image" validate:"required"`
	Email    string `form:"email" json: "email" validate:"required"`
	Role     string `form:"role" json: "role" validate:"required"`
	Password string `form:"password" json: "password" validate:"required"`
	//default false (Active)
	Disable  bool `gorm:"default:0"`
	TenantID uint
}

func CreateUser(db *gorm.DB, newUser *User) (err error) {
	err = db.Create(newUser).Error
	if err != nil {
		return err
	}
	return nil
}

func FindUserByUsername(db *gorm.DB, user *User, username string) (err error) {
	err = db.Where("username=?", username).First(user).Error
	if err != nil {
		return err
	}
	return nil
}

// update user data with username
func UpdateUserByUsername(db *gorm.DB, user *User, username string) (err error) {
	db.Where("username=?", username).First(&user)
	err = db.Save(user).Error
	if err != nil {
		return err
	}
	return nil
}
