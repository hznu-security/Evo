package model

import "gorm.io/gorm"

//const TYPE_AWD int = 0
//const TYPE_SOLVE int = 1

type Challenge struct {
	gorm.Model
	Title       string  `gorm:"type:varchar(100)" binding:"required,max=100"`
	Desc        string  `gorm:"type:varchar(255)" binding:"required,max=255"` // 题目描述
	AutoRefresh bool    //是否自动刷新flag
	Command     string  `gorm:"type:varchar(255)" binding:"max=255"` //刷新flag时使用的shell命令
	Visible     bool    //是否可见
	Score       float64 `binding:"required"`
}

func (Challenge) TableName() string {
	return "challenges"
}
