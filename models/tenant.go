package models

import (
	"gorm.io/gorm"
)

type Tenant struct {
	gorm.Model
	Id    int    `form:"id" json: "id" validate:"required"`
	Name  string `form:"name" json: "name" validate:"required"`
	Users []User
}
