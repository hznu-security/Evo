package model

import "gorm.io/gorm"

// GameBox 一个或多个gamebox对应一个实际的容器
type GameBox struct {
	gorm.Model
	ChallengeID uint
	TeamId      uint
	CName       string `gorm:"type:varchar(100)"` // 其对应容器的name
	Name        string `gorm:"type:varchar(100)"`
	Port        string `gorm:"type:varchar(100)"`
	SshPort     string `gorm:"type:varchar(20)"`
	SshUser     string `gorm:"type:varchar(20)"`
	SshPwd      string `gorm:"type:varchar(50)"`
	Score       float64
	Visible     bool
	IsDown      bool
	IsAttacked  bool
}

func (GameBox) TableName() string {
	return "game_boxes"
}
