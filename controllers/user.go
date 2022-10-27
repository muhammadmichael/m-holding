package controllers

import (
	"fmt"
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

	params := c.AllParams() // "{"username": "admin"}"
	intId, _ := strconv.Atoi(params["id"])
	user.Id = intId
	// username := params["username"]

	if err := c.BodyParser(&user); err != nil {
		// send error message
		return c.JSON(fiber.Map{
			"message": "error",
		})
	}
	Name := ""
	Email := ""
	Role := ""
	Username := ""
	// Image := ""
	// Parse the multipart form:
	if form, err := c.MultipartForm(); err == nil {
		// contains non-file fields
		// fill user value with form value
		Name = form.Value["name"][0]
		Email = form.Value["email"][0]
		Role = form.Value["role"][0]
		Username = form.Value["username"][0]
		files := form.File["image"]
		for _, file := range files {
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			// Save files to disk:
			user.Image = fmt.Sprintf("%s", file.Filename)

			if err := c.SaveFile(file, "./public/uploads/profilepict/"+file.Filename); err != nil {
				return err
			}
			//find user by id
			err := models.FindUserById(controller.Db, &user, intId)
			if err != nil {
				return c.SendStatus(500) // http 500 internal server error
			}
			user.Name = Name
			user.Email = Email
			user.Role = Role
			user.Username = Username
			user.Password = user.Password
			user.Image = fmt.Sprintf("%s", file.Filename)
			//bcrypt password

			// save user
			errs := models.UpdateUser(controller.Db, &user)
			if errs != nil {
				return c.JSON(fiber.Map{
					"message": "error",
				})
			}
		}
	}
	// if succeed
	return c.Render("profile", fiber.Map{
		"Title": "Profile",
		"User":  user,
	})
}

// delete user
func (controller *UserController) DeleteUser(c *fiber.Ctx) error {
	params := c.AllParams() // "{"id": "1"}"

	intId, _ := strconv.Atoi(params["id"])

	var user models.User
	err := models.FindUserById(controller.Db, &user, intId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	errs := models.DeleteUser(controller.Db, &user, intId)
	if errs != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	return c.Redirect("/users")
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

// update profile picture user
func (u *UserController) UpdateProfilePictureUser(id uint, image string) (err error) {

	err = u.Db.Model(&models.User{}).Where("id=?", id).Update("image", image).Error //update image
	if err != nil {
		return err
	}
	return nil
}

// update user data & profile picture with username
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
func (u *UserController) DisableTheUser(c *fiber.Ctx) (err error) {
	id, _ := strconv.Atoi(c.Params("id"))
	fmt.Println(id)
	err = u.Db.Model(&models.User{}).Where("id=?", id).Update("disable", true).Error //update disable
	if err != nil {
		return err
	}
	return c.Redirect("/login")
}

// enable user
func (u *UserController) EnableTheUser(c *fiber.Ctx) (err error) {
	id, _ := strconv.Atoi(c.Params("id"))
	err = u.Db.Model(&models.User{}).Where("id=?", id).Update("disable", false).Error //update disable
	if err != nil {
		return err
	}
	return c.Redirect("/login")
}
