package api

import (
	"rapid/m-holding/database"
	"rapid/m-holding/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


type UserApiController struct {
	// Declare variables
	Db *gorm.DB
}


func InitUserApiController() *UserApiController {
	db := database.InitDb()
	// gorm sync
	db.AutoMigrate(&models.User{})

	return &UserApiController{Db: db}
}

func (controller *UserApiController) RegisterApi(c *fiber.Ctx) error {
	//myform := new(models.Product)
	var registrasi models.User

	if err := c.BodyParser(&registrasi); err != nil {
		return c.JSON(fiber.Map{
			"Message" : "Failed",
			"Status"   :500,	

		})
	}
	// save registrasi
	err := models.Registrasi(controller.Db, &registrasi)
	if err != nil {
		return c.JSON(fiber.Map{
			"Message" : "Failed",
			"Status"   :500,	

		})
	}
	// if succeed
	return c.JSON(fiber.Map{
		"Message" : "Success",
		"Data" : registrasi,
		"Status"   :200,	

	})
}