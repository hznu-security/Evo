package model

import "time"

type Flag struct {
	ID          uint `gorm:"primarykey"`
	TeamId      uint
	BoxId       uint
	ChallengeID uint
	Round       uint
	Flag        string `gorm:"varchar(255);not null"`
	DeletedAt   time.Time
}
