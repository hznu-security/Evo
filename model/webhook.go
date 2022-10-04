package model

import "gorm.io/gorm"

type Webhook struct {
	gorm.Model
	Url     string  `gorm:"type:varchar(255)" binding:"required,max=255"`
	Type    string  `gorm:"type:varchar(30)" binding:"required,max=30"`
	Time    uint    `binding:"required"` //重试次数
	Timeout float64 `binding:"required"` //时限 秒
}
