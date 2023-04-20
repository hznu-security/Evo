package model

import "gorm.io/gorm"

type Attack struct {
	gorm.Model
	TeamID      uint `gorm:"type:varchar(50);not null"`
	Attacker    uint
	GameBoxId   uint
	ChallengeId uint
	Round       uint
}

func (Attack) TableName() string {
	return "attacks"
}
