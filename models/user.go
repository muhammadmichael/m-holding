package models

import (
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	"fmt"
)

type User struct {
	gorm.Model
	Name     string `form:"name" json: "name" validate:"required"`
	Username string `form:"username" json: "username" validate:"required"`
	Image    string `form:"image" json: "image" validate:"required"`
	Email    string `form:"email" json: "email" validate:"required"`
	Role     string `form:"role" json: "role" validate:"required"`
	Password string `form:"password" json: "password" validate:"required"`
	Disabel  bool `gorm:"default:0"`
	TenantID uint
}

func Registrasi(db *gorm.DB, newUser *User) (err error) {
	plainPassword := newUser.Password
	bytes, _ := bcrypt.GenerateFromPassword([]byte(plainPassword),10)
	sHash := string(bytes)
	fmt.Println("Hash password: ", sHash)
	newUser.Password = sHash
	err = db.Create(newUser).Error
	if err != nil {
		return err
	}
	return nil
}

// func FindUserByUsername(db *gorm.DB, user *User, username string) (err error) {
// 	err = db.Where("username=?", username).First(user).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
