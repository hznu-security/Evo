package model

type Score struct {
	ID        uint `gorm:"primarykey"`
	TeamId    uint
	GameBoxId uint
	Round     uint
	Score     float64
	Reason    string `gorm:"type:varchar(30)"`
}
