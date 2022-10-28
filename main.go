package main

import (
	"fmt"
	"rapid/m-holding/controllers"

	"github.com/gofiber/fiber/v2"
	//"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"
)

type Product struct {
	Id      int
	Name    string
	Viewer  int
	Revenue float32
}

func main() {
	// session
	//store := session.New()

	// load template engine
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// static
	app.Static("/", "./public", fiber.Static{
		Index: "",
	})

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
	//userController := controllers.InitUserController(store)
	tenantController := controllers.InitTenantController()
	//userApiController := api.InitUserApiController()

	// Test
	fmt.Println(tenantController)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("home", fiber.Map{
			"Title": "M-Holding",
		})
	})

	//
	DashboardController := controllers.InitDashboardController()
	app.Get("/dashboard", DashboardController.Dashboard)

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.Render("test", fiber.Map{
			"Title": "Revenue",
		})
	})

	app.Get("/revenue", func(c *fiber.Ctx) error {
		var products = []*Product{
			{Id: 1, Name: "Iklan 1", Viewer: 10, Revenue: 20},
			{Id: 2, Name: "Iklan 2", Viewer: 20, Revenue: 40},
			{Id: 3, Name: "Iklan 3", Viewer: 30, Revenue: 60},
			{Id: 4, Name: "Iklan 4", Viewer: 40, Revenue: 80},
			{Id: 5, Name: "Iklan 5", Viewer: 50, Revenue: 100},
		}

		return c.Render("revenue", fiber.Map{
			"Title":    "Detail Revenue",
			"Products": products,
		})
	})

	app.Listen(":3000")
}
