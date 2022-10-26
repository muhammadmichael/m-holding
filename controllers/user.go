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

// update profile picture user
func (u *UserController) UpdateProfilePictureUser(id uint, image string) (err error) {

	err = u.Db.Model(&models.User{}).Where("id=?", id).Update("image", image).Error //update image
	if err != nil {
		return err
	}
	return nil
}

// update user profile picture with username
func (controller *UserController) UpdateUserImage(c *fiber.Ctx) error {
	var user models.User

	params := c.AllParams() // "{"username": "admin"}"
	username := params["username"]

	if err := c.BodyParser(&user); err != nil {
		// send error message
		return c.JSON(fiber.Map{
			"message": "error",
		})
	}

	// Parse the multipart form:
	if form, err := c.MultipartForm(); err == nil {
		// => *multipart.Form
		// Get all files from "image" key:
		files := form.File["image"]
		// => []*multipart.FileHeader
		// Loop through files:
		for _, file := range files {
			// Save files to ./uploads folder
			// => *multipart.FileHeader
			// => error
			if err := c.SaveFile(file, "./uploads/profilepict"+file.Filename); err != nil {
				return err
			}
		}
	}
	// save user
	err := models.UpdateUserByUsername(controller.Db, &user, username)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "error",
		})
	}
	// if succeed
	return c.JSON(fiber.Map{
		"message": "success update profile pict",
	})
}

// disable user
func (u *UserController) DisableUser(id uint) (err error) {

	err = u.Db.Model(&models.User{}).Where("id=?", id).Update("disable", true).Error //update disable
	if err != nil {
		return err
	}
	return nil
}
