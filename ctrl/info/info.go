package info

import (
	"Evo/config"
	"Evo/db"
	"Evo/model"
	"Evo/util"
	"github.com/gin-gonic/gin"
	"sort"
)

// Time 返回比赛时间信息
/**
开始时间
结束时间
当前第几轮
本轮剩余
*/
func Time(c *gin.Context) {
	type timeInfo struct {
		BeginTime       int64
		EndTime         int64
		NowRound        uint
		RoundTime       uint    // 每轮时间，单位秒
		RoundRemainTime float64 // 单位为秒
	}
	util.Success(c, "success", gin.H{
		"time": timeInfo{
			BeginTime:       config.StartTime.Unix(),
			EndTime:         config.EndTime.Unix(),
			NowRound:        config.ROUND_NOW,
			RoundTime:       config.ROUND_TIME,
			RoundRemainTime: config.GetRoundRemainTime(),
		},
	})
}

type Team struct {
	Id    int
	Name  string
	Rank  int
	Img   string // 队伍logo 的url
	Score int
}

// GetRank 获取排名 // TODO
func GetRank(c *gin.Context) {
	var teams []Team
	var team []model.Team
	db.DB.Select([]string{"id", "name", "logo", "score"}).Find(&team)
	// 根据分数排序后返回
	sort.Slice(team, func(i, j int) bool {
		return team[i].Score > team[j].Score
	})

	for rank, team := range team {
		teams = append(teams, Team{
			Id:    int(team.ID),
			Name:  team.Name,
			Rank:  rank + 1,
			Img:   team.Logo,
			Score: int(team.Score),
		})
	}

	util.Success(c, "success", gin.H{
		"rank": teams,
	})
}
