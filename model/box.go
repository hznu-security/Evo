package model

import "gorm.io/gorm"

// Box 每个box对应一个实际意义上的容器
type Box struct {
	gorm.Model
	ChallengeID uint
	TeamId      uint
	Name        string `gorm:"type:varchar(100)"`
	Port        string `gorm:"type:varchar(100)"` // 题目的port
	SshPort     string `gorm:"type:varchar(20)"`  // ssh登录的port
	SshUser     string `gorm:"type:varchar(20)"`
	SshPwd      string `gorm:"type:varchar(50)"`
}

func (Box) TableName() string {
	return "boxes"
}
