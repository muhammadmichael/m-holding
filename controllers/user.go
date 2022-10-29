package controllers

import (
	"fmt"
	"rapid/m-holding/database"
	"rapid/m-holding/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
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
	var user models.User

	if err := c.BodyParser(&registrasi); err != nil {
		return c.Redirect("/auth")
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
		return c.Redirect("/auth")
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
		sess.Set("userId", user.Id)
		sess.Save()

		return c.Redirect("/")
	}

	return c.Redirect("/login")
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

// GET /
// Home
func (controller *UserController) GetHome(c *fiber.Ctx) error {
	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}
	val := sess.Get("userId")

	return c.Render("home", fiber.Map{
		"Title":  "M-Holding",
		"UserId": val,
	})
}

func (controller *UserController) DataUser(c *fiber.Ctx) error {
	var user []models.User

	err := models.ReadUser(controller.Db, &user)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.Render("indexuser", fiber.Map{
		"Title": "Data User",
		"User":  user,
	})
}

func (controller *UserController) UserDisable(c *fiber.Ctx) (err error) {
	id, _ := strconv.Atoi(c.Params("id"))

	err = controller.Db.Model(&models.User{}).Where("id=?", id).Update("disable", true).Error //update disable
	if err != nil {
		return err
	}
	return c.Redirect("/user")
}

// enable user
func (controller *UserController) UserEnable(c *fiber.Ctx) (err error) {
	id, _ := strconv.Atoi(c.Params("id"))
	err = controller.Db.Model(&models.User{}).Where("id=?", id).Update("disable", false).Error //update disable
	if err != nil {
		return err
	}
	return c.Redirect("/user")
}

func (controller *UserController) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)

	var user models.User
	models.DeleteById(controller.Db, &user, idn)

	//return c.JSON(user)
	return c.Redirect("/user")
}

// GET FORM REGISTRASI
func (controller *UserController) UpdateUserForm(c *fiber.Ctx) error {

	id := c.Params("id")
	idn, _ := strconv.Atoi(id)

	var user models.User
	err := models.FindUserById(controller.Db, &user, idn)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.Render("Update", fiber.Map{
		"Title": "Update User",
		"User":  user,
	})
}

func (controller *UserController) EditUser(c *fiber.Ctx) error {
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)

	var user models.User
	err := models.FindUserById(controller.Db, &user, idn)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	var updateUser models.User

	if err := c.BodyParser(&updateUser); err != nil {
		return c.SendStatus(400)
	}
	user.Name = updateUser.Name
	user.Username = updateUser.Username
	user.Email = updateUser.Email
	user.Role = updateUser.Role
	user.KategoriUser = updateUser.KategoriUser

	// save suer
	models.UpdateUser(controller.Db, &user)

	return c.Redirect("/user")
}
