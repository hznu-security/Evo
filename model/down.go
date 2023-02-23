package model

type Down struct {
	ID        uint
	GameBoxId uint
}

func (Down) TableName() string {
	return "downs"
}
