package api

import (
	"rapid/m-holding/database"
	"rapid/m-holding/models"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginForm struct {
	Username string `form:"username" json:"username" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
}

type UserApiController struct {
	// Declare variables
	Db *gorm.DB
}

var checker = validator.New()

func InitUserApiController() *UserApiController {
	db := database.InitDb()
	// gorm sync
	db.AutoMigrate(&models.User{})

	return &UserApiController{Db: db}
}

func (controller *UserApiController) RegisterApi(c *fiber.Ctx) error {
	//myform := new(models.Product)
	var registrasi models.User
	var user models.User

	if err := c.BodyParser(&registrasi); err != nil {
		return c.JSON(fiber.Map{
			"Message": "Failed",
			"Status":  500,
		})
	}

	// Cek apakah username sudah digunakan
	errUsername := models.FindUserByUsername(controller.Db, &user, registrasi.Username)
	if errUsername != gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"message": "Username telah digunakan",
		})
	}

	// save registrasi
	err := models.Registrasi(controller.Db, &registrasi)
	if err != nil {
		return c.JSON(fiber.Map{
			"Message": "Failed",
			"Status":  500,
		})
	}
	// if succeed
	return c.JSON(fiber.Map{
		"Message": "Success",
		"Data":    registrasi,
		"Status":  200,
	})
}

// POST /api/login
func (controller *UserApiController) ApiLoginPosted(c *fiber.Ctx) error {
	var user models.User
	var myform LoginForm

	if err := c.BodyParser(&myform); err != nil {
		// Bad Request, LoginForm is not complete
		return c.JSON(fiber.Map{
			"status":  400,
			"message": "Bad Request, Login Form is not complete",
		})
	}

	err := checker.Struct(myform)
	if err != nil {
		return c.JSON(fiber.Map{
			"status":  400,
			"message": "Bad Request, Login Form is not complete",
		})
	}

	// Find user
	errs := models.FindUserByUsername(controller.Db, &user, myform.Username)
	if errs != nil {
		return c.JSON(fiber.Map{
			"message": "Cannot find user",
		})
	}

	// Compare password
	compare := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(myform.Password))
	if compare == nil { // compare == nil artinya hasil compare di atas true
		// Create the Claims
		exp := time.Now().Add(time.Hour * 72) // token expired time: 72 hours
		claims := jwt.MapClaims{
			"name":  user.Username,
			"admin": true,
			"exp":   exp.Unix(),
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("mysecretpassword"))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{
			"message": "Berhasil Login",
			"token":   t,
			"expired": exp.Format("2006-01-02 15:04:05"),
		})
	}

	return c.JSON(fiber.Map{
		"status":  401,
		"message": "Unauthorized",
	})
}
