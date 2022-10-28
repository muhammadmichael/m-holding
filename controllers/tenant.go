package controllers

import (
	"rapid/m-holding/database"
	"rapid/m-holding/models"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type TenantController struct {
	// Declare variables
	Db *gorm.DB
}

var checker = validator.New()

func InitTenantController() *TenantController {
	db := database.InitDb()
	// gorm sync
	db.AutoMigrate(&models.Tenant{})

	return &TenantController{Db: db}
}

