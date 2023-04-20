package model

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Name  string `gorm:"type:varchar(30);not null" binding:"required,max=30"`
	Pwd   string `gorm:"type:varchar(255);not null" binding:"required,max=30"`
	Token string `gorm:"type:varchar(255)"`
}

func (Admin) TableName() string {
	return "admins"
}
