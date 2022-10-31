package controllers

import (
	//"rapid/m-holding/database"
	//"rapid/m-holding/models"

	//"github.com/go-playground/validator/v10"
	"rapid/m-holding/database"
	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

type DashboardController struct {
	// Declare variables
	Db *gorm.DB
}

//var checker = validator.New()

func InitDashboardController() *DashboardController {
	db := database.InitDb()
	// gorm sync
	//db.AutoMigrate(&models.Dashboard{})

	return &DashboardController{Db: db}
}

func (controller *DashboardController) Dashboard(c *fiber.Ctx) error {
	return c.Render("dashboard", fiber.Map{
		"Title": "Dashboard",
	})
}