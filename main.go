package main

import (
	"fmt"
	"rapid/m-holding/controllers"
	"rapid/m-holding/api"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {
	// session
	// store := session.New()

	// load template engine
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// static
	app.Static("/", "./public", fiber.Static{
		Index: "",
	})

	// controllers
	userController := controllers.InitUserController()
	tenantController := controllers.InitTenantController()
	UserApiController := api.InitUserApiController()

	fmt.Println(userController)
	fmt.Println(tenantController)

	//Module User Api
	api:= app.Group("/api")
	api.Post("/register", UserApiController.RegisterApi)
	

	//Module User
	user := app.Group("/register")
	user.Get("/", userController.Register)
	user.Post("/tambah", userController.NewRegister)

	login := app.Group("/login")
	login.Get("/", userController.Login)
	//login.Post("/sigIn", userController.LoginUser)

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.Render("home", fiber.Map{
			"Title": "M-Holding",
		})
	})

	app.Listen(":3000")
}
