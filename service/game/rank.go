package game

import (
	"Evo/db"
	"Evo/model"
)

type RankItem struct {
	Id    uint
	Name  string
	Rank  int
	Img   string
	Score float64
}

func SetRankList() {
	var rankList []*RankItem
	var teams []model.Team
	db.DB.Model(&model.Team{}).Order("score DESC").Find(&teams)
	for i, team := range teams {
		rankList = append(rankList, &RankItem{
			Id:    team.ID,
			Name:  team.Name,
			Rank:  i + 1,
			Score: team.Score,
		})
	}
	db.Set("rankList", rankList)
}

func GetRankList() []*RankItem {
	rankList, ok := db.Get("rankList")
	if !ok {
		SetRankList()
		rankList, _ = db.Get("rankList")
	}
	return rankList.([]*RankItem)
}
