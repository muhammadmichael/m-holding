package controllers

import (
	"rapid/m-holding/database"
	"rapid/m-holding/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	// "fmt"
  // "reflect"
)

type TenantForm struct {
	Name string `form:"name" json:"name" validate:"required"`
	Users string `form:"user" json:"user" validate:"required"`
}

type TenantController struct {
	// Declare variables
	Db *gorm.DB
	// store *session.Store

	Id			int
	Name		string
	User		string

}

// var checker = validator.New()

func InitTenantController() *TenantController {
	db := database.InitDb()
	// // gorm sync
	db.AutoMigrate(&models.Tenant{})

	return &TenantController{Db: db}
}

//GET AllTenant
func (controller *TenantController) AllTenant(c *fiber.Ctx) error {
	var tenants []models.Tenant
	err := models.ReadAllTenant(controller.Db, &tenants)
	if err!=nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	// var tenant = []*TenantController{
	// 	{Id: 1, Name: "ven", User: "asdas,dfgdf,awdaw,sdgsd"},
	// 	{Id: 2, Name: "ven 2", User: "asdas,dfgdf,awdaw,sdgsd"},
	// 	{Id: 3, Name: "ven 3", User: "asdas,dfgdf,awdaw,sdgsd"},
	// 	{Id: 4, Name: "ven 4", User: "asdas,dfgdf,awdaw,sdgsd"},
	// }

	return c.Render("indextenant", fiber.Map{ // abaikan alamt view 
		"Title": "M-Holding",
		"Tenants": tenants,
	})

	// API
	// return c.JSON(fiber.Map{
	// 	"Message":  "Berhasil mendapatkan seluruh list products",
	// 	"Tenant": tenant,
	// })
}

//GET AddTenant
func (controller *TenantController) AddTenant(c *fiber.Ctx) error {
	return c.Render("addtenant", fiber.Map{ // abaikan view 
		"Title": "M-Holding",
	})
}

//POST AddTenant 
func (controller *TenantController) AddTenantPosted(c *fiber.Ctx) error {
	var myform models.Tenant

	if err := c.BodyParser(&myform); err != nil {
		return c.Redirect("/tenant") //abaikan view
		// API
		// return c.JSON(fiber.Map{
		// 	"status":  400,
		// 	"message": "Bad Request, Tenant Form is not complete",
		// })
	}

	errr := models.CreateTenant(controller.Db, &myform)
	if errr != nil{
		return c.Redirect("/tenant") // abaikan view
		// API
		// return c.SendStatus(500)
		
	}
	return c.Redirect("/tenant") // abaikan view
	// API 
	// return c.JSON(fiber.Map{
	// 	"status":  200,
	// 	"message": "Berhasil Menambahkan Product",
	// })
}

// GET Tenant by Id
func (controller *TenantController) DetailTenant(c *fiber.Ctx)error{
	id := c.Query("id")
	idn,_ := strconv.Atoi(id)

	var tenant models.Tenant
	err := models.FindTenantById(controller.Db, &tenant, idn)
	if err != nil {
		return c.SendStatus(500)
		// API
		// return c.JSON(fiber.Map{
		// 	"Status":  500,
		// 	"message": "Tidak ditemukan tenant dengan Id" + id,
		// }) 
	}
	return c.Render("tenantdetailid", fiber.Map{ // abaikan view
		"Title": "M-Holding",
		"Tenant": tenant,
	})
	// API 
	// return c.JSON(fiber.Map{
	// 	"message": "Detail tenant dengan Id " + id,
	// 	"Tenant": tenant,
	// })
}

// GET Tenant by name
func (controller *TenantController) DetailTenant2(c *fiber.Ctx)error{
	name := c.Query("name")

	var tenant models.Tenant
	err := models.FindTenantByName(controller.Db, &tenant, name)
	if err != nil {
		return c.SendStatus(500)
		// API
		// return c.JSON(fiber.Map{
		// 	"Status":  500,
		// 	"message": "Tidak ditemukan tenant dengan name" + name,
		// }) 
	}
	return c.Render("tenantdetailname", fiber.Map{ // abaikan view
		"Title": "M-Holding",
		"Tenant": tenant,
	})
	// API 
	// return c.JSON(fiber.Map{
	// 	"message": "Detail tenant dengan name " + name,
	// 	"Tenant": tenant,
	// })
}

//POST edit tenant
func (controller *TenantController) EditTenantPosted(c *fiber.Ctx)error{
	id := c.Params("id")
	idn,_ := strconv.Atoi(id)


	var tenant models.Tenant
	err := models.FindTenantById(controller.Db, &tenant, idn)
	if err!=nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	var myform models.Tenant

	if err := c.BodyParser(&myform); err != nil {
		return c.Redirect("/tenants") //abaikan view
	}
	tenant.Name = myform.Name
	tenant.Users = myform.Users
	// save product
	models.AddUserTenant(controller.Db, &tenant)
	
	return c.Redirect("/tenants")	
}