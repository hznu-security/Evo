package model

import "gorm.io/gorm"

type Team struct {
	gorm.Model
	Name  string  `gorm:"type:varchar(50);not null"  binding:"required,max=50"`
	Pwd   string  `gorm:"type:varchar(255);not null"`
	Logo  string  `gorm:"type:varchar(50)"` //logo 的地址
	Score float64 `gorm:"not null"`
	Token string  `gorm:"type:text;not null"`
}

func (Team) TableName() string {
	return "teams"
}
