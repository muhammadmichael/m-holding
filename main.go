package main

import (
	"rapid/m-holding/api"
	"rapid/m-holding/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
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
	dashboardController := controllers.InitDashboardController()
	userApiController := api.InitUserApiController()

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
	profile.Delete("/delete/:id", userController.DeleteUser)
	//disable user
	profile.Get("/disable/:id", userController.UserDisable)
	profile.Get("/enable/:id", userController.UserEnable)

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

	tnt := app.Group("/tenants")
	tnt.Get("/", tenantController.AllTenant)

	dashboard := app.Group("dashboard")
	dashboard.Get("/", dashboardController.Dashboard)

	app.Listen(":3000")
}
