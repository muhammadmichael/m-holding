package main

import (
	"fmt"
	"rapid/m-holding/controllers"

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

	fmt.Println(userController)
	fmt.Println(tenantController)

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.Render("home", fiber.Map{
			"Title": "M-Holding",
		})
	})

	// Auth Routes (Register and Login)
	app.Post("/login", userController.LoginPosted)
	app.Post("/register", userController.AddRegisteredUser)

	app.Listen(":3000")
}
