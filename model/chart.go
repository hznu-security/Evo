package model

type Chart struct {
	ID       uint   `gorm:"primarykey"`
	TeamName string `gorm:"type:varchar(50)"`
	Round    uint
	Score    float64
}

func (Chart) TableName() string {
	return "charts"
}
