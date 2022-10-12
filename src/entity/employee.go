package entity

import (
	"time"
)

func (Employee) TableName() string {
	return "Employee"
}

type Employee struct {
	Username string `gorm:"primaryKey"`
	Password string
	FirstName string
	LastName string
	RoleID int
	ReportsTo string
	LastLogin time.Time
	CreatedAt time.Time
	DeletedAt time.Time
}

func (Employee) GetAll() []Employee {
	var employees []Employee
	db.Find(&employees)

	return employees
}