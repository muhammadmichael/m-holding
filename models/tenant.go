package models

import (
	// "fmt"

	// "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Tenant struct {
	gorm.Model
	Id    int    `form:"id" json:"id" validate:"required"`
	Name  string `form:"name" json:"name" validate:"required"`
	// e.g. : ali, ili, uli, oli 
	Users string `forn:"user" json:"name"`
}

func ReadAllTenant(db *gorm.DB, tenant *[]Tenant)(err error){
	err = db.Find(tenant).Error
	if err != nil {
		return err
	}
	return nil
}

func FindTenantById(db *gorm.DB, tenant *Tenant, id int) (err error) {
	err = db.Where("id=?", id).First(tenant).Error
	if err != nil {
		return err
	}
	return nil
}

// func UpdateTenant(db *gorm.DB, tenant *Tenant) (err error) {
// 	db.Save(tenant)
// 	return nil
// }

func FindTenantByName(db *gorm.DB, tenant *Tenant, name string) (err error) {
	err = db.Where("name=?", name).First(tenant).Error
	if err != nil {
		return err
	}
	return nil
}

func CreateTenant(db *gorm.DB, newTenant *Tenant)(err error){
	err = db.Create(newTenant).Error
	if err != nil {
		return err
	}
	return nil
}
//menggunakan update
func AddUserTenant(db *gorm.DB, tenant *Tenant) (err error) {
	db.Save(tenant)
	return nil
}

func DeleteTenantById(db *gorm.DB, tenant *Tenant ,id int)(err error){
	db.Where("id=?",id).Delete(tenant)
	return nil
}