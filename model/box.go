package model

import "gorm.io/gorm"

type Box struct {
	gorm.Model
	ChallengeID uint
	TeamId      uint
	Name        string `gorm:"type:varchar(100)"`
	Ip          string `gorm:"type:varchar(20)"`
	Port        string `gorm:"type:varchar(20)"`
	SshUser     string `gorm:"type:varchar(20)"`
	SshPwd      string `gorm:"type:varchar(50)"`
	Score       float64
	IsDown      bool
	IsAttacked  bool
}
