package controllers

import (
	"rapid/m-holding/database"
	"rapid/m-holding/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginForm struct {
	Username string `form:"username" json:"username" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
}

type UserController struct {
	// Declare variables
	Db    *gorm.DB
	store *session.Store
}

func InitUserController(s *session.Store) *UserController {
	db := database.InitDb()
	// gorm sync
	db.AutoMigrate(&models.User{})

	return &UserController{Db: db, store: s}
}

// GET FORM REGISTRASI
func (controller *UserController) Register(c *fiber.Ctx) error {
	return c.Render("registrasi", fiber.Map{
		"Title": "Register User",
	})
}

// POST TO REGISTRASI
func (controller *UserController) NewRegister(c *fiber.Ctx) error {
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

// GET FORM LOGIN /login
func (controller *UserController) Login(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"Title": "Login",
	})
}

// POST /login
func (controller *UserController) LoginPosted(c *fiber.Ctx) error {
	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}
	var user models.User
	var myform LoginForm

	if err := c.BodyParser(&myform); err != nil {
		return c.Redirect("/login")
	}

	// Find user
	errs := models.FindUserByUsername(controller.Db, &user, myform.Username)
	if errs != nil {
		return c.Redirect("/login") // Unsuccessful login (cannot find user)
	}

	// Compare password
	compare := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(myform.Password))
	if compare == nil { // compare == nil artinya hasil compare di atas true
		sess.Set("username", user.Username)
		sess.Set("userId", user.ID)
		sess.Save()

		return c.Redirect("/")
	}

	return c.Redirect("/login")
}

// POST /api/login
func (controller *UserController) ApiLoginPosted(c *fiber.Ctx) error {
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

// GET /profile/:id
func (controller *UserController) ViewProfile(c *fiber.Ctx) error {
	params := c.AllParams() // "{"id": "1"}"

	intId, _ := strconv.Atoi(params["id"])

	var user models.User
	err := models.FindUserById(controller.Db, &user, intId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}
	val := sess.Get("userId")

	return c.Render("profile", fiber.Map{
		"Title":  "Profile",
		"User":   user,
		"UserId": val,
	})
}

// GET /profile/edit/:id
func (controller *UserController) EditProfile(c *fiber.Ctx) error {
	params := c.AllParams() // "{"id": "1"}"

	intId, _ := strconv.Atoi(params["id"])

	var user models.User
	err := models.FindUserById(controller.Db, &user, intId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}
	val := sess.Get("userId")

	return c.Render("editprofile", fiber.Map{
		"Title":  "Edit Profile",
		"User":   user,
		"UserId": val,
	})
}

// POST /profile/edit/:id
func (controller *UserController) EditProfilePosted(c *fiber.Ctx) error {
	var user models.User
	var userEditForm models.User

	params := c.AllParams() // "{"id": "1"}"
	intId, _ := strconv.Atoi(params["id"])
	user.Id = intId

	if err := c.BodyParser(&userEditForm); err != nil {
		return c.Redirect("/profile/edit/" + params["id"])
	}

	// Find the user
	err := models.FindUserById(controller.Db, &user, intId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	// Change from user's input
	user.Name = userEditForm.Name
	user.Email = userEditForm.Email
	user.Role = userEditForm.Role

	// save product
	errs := models.UpdateUser(controller.Db, &user)
	if errs != nil {
		return c.Redirect("/profile/edit/" + params["id"])
	}

	// if succeed
	return c.Redirect("/profile/" + params["id"])
}

// /logout
func (controller *UserController) Logout(c *fiber.Ctx) error {

	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}
	sess.Destroy()
	return c.Render("login", fiber.Map{
		"Title": "Login",
	})
}
