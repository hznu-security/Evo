package model

type Chart struct {
	ID       uint   `gorm:"primarykey"`
	TeamName string `gorm:`
}

func (Chart) TableName() string {
	return "charts"
}
