package model

type Attack struct {
	ID          uint `gorm:"primarykey"`
	TeamID      uint `gorm:"type:varchar(50);not null"`
	Attacker    uint
	BoxId       uint
	ChallengeId uint
	Round       uint
}

func (Attack) TableName() string {
	return "attacks"
}
