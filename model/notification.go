package model

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	Title   string `gorm:"varchar(100);not null" binding:"required"`
	Content string `gorm:"varchar(255);not null" binding:"required"`
}

func (Notification) TableName() string {
	return "notifications"
}
