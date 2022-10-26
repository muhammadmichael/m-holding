package controllers

import (
	"rapid/m-holding/database"
	"rapid/m-holding/models"
	"github.com/gofiber/fiber/v2"
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

//GET FORM REGISTRASI
func (controller *UserController) Register(c *fiber.Ctx) error {
	return c.Render("registrasi", fiber.Map{
		"Title": "Register User",
	})
}

//POST TO REGISTRASI
func (controller *UserController) NewRegister(c *fiber.Ctx) error {
	//myform := new(models.Product)
	var registrasi models.User

	if err := c.BodyParser(&registrasi); err != nil {
		return c.Redirect("/registrasi")
	}
	// save registrasi
	err := models.Registrasi(controller.Db, &registrasi)
	if err != nil {
		return c.Redirect("/registrasi")
	}
	// if succeed
	return c.Redirect("/login")
}

func (controller *UserController) Login(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"Title": "Login",
	})
}

