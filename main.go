package main

import (
	"fmt"
	"rapid/m-holding/api"
	"rapid/m-holding/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"
)

func main() {
	// session
	store := session.New()

	// load template engine
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// static
	app.Static("/public", "./public")

	// Middleware to check login
	// CheckLogin := func(c *fiber.Ctx) error {
	// 	sess, _ := store.Get(c)
	// 	val := sess.Get("username")
	// 	if val != nil {
	// 		return c.Next()
	// 	}

	// 	return c.Redirect("/login")
	// }

	// controllers
	userController := controllers.InitUserController(store)
	tenantController := controllers.InitTenantController()
	userApiController := api.InitUserApiController()

	// Test
	fmt.Println(tenantController)

	// Home
	app.Get("/", userController.GetHome)

	// Auth Routes (Register and Login)
	app.Get("/register", userController.Register)
	app.Post("register/tambah", userController.NewRegister)
	app.Get("/login", userController.Login)
	app.Get("/logout", userController.Logout)
	app.Post("/login", userController.LoginPosted)

	profile := app.Group("/profile")
	profile.Get("/:id", userController.ViewProfile)
	profile.Get("/edit/:id", userController.EditProfile)
	profile.Post("/edit/:id", userController.EditProfilePosted)

	// API Routes
	api := app.Group("/api")
	api.Post("/login", userApiController.ApiLoginPosted)
	api.Post("/register", userApiController.RegisterApi)

	usercont := app.Group("/user")
	usercont.Get("/", userController.DataUser)
	usercont.Get("/enable/:id", userController.UserDisable)
	usercont.Get("/disable/:id", userController.UserEnable)
	usercont.Get("/deleteuser/:id", userController.DeleteUser)
	usercont.Get("/update/:id", userController.UpdateUserForm)
	usercont.Post("/edituser/:id", userController.EditUser)

	app.Listen(":3000")
}
