package controllers

import (
	"rapid/m-holding/database"
	"rapid/m-holding/models"

	"gorm.io/gorm"
)

type UserController struct {
	// Declare variables
	Db *gorm.DB
}

func InitUserController() *UserController {
	db := database.InitDb()
	// gorm sync
	db.AutoMigrate(&models.User{})

	return &UserController{Db: db}
}
