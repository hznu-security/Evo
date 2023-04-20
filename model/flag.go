package model

type Flag struct {
	ID          uint `gorm:"primarykey"`
	TeamId      uint
	GameBoxId   uint
	ChallengeID uint
	Round       uint
	Flag        string `gorm:"type:varchar(255);not null"`
}

func (Flag) TableName() string {
	return "flags"
}
