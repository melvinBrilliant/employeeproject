package entity

import (
	"com.melvin.employee/src/util"
)

var db = util.ConnectToDatabase()

func (Role) TableName() string {
	return "Role"
}

type Role struct {
	ID int `gorm:"column:ID"`
	RoleName string 
	Description string
}

func (Role) GetAll() []Role {
	var roles []Role
	db.Find(&roles)

	return roles
}