package model

import "gorm.io/gorm"

type Down struct {
	gorm.Model
	GameBoxId   uint
	Round       uint
	TeamId      uint
	ChallengeId uint
}

func (Down) TableName() string {
	return "downs"
}
