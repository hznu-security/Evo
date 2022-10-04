package model

type Attack struct {
	ID          uint `gorm:"primarykey"`
	TeamID      uint
	Attacker    uint //攻击者
	BoxId       uint
	ChallengeId uint
	Round       uint
}
