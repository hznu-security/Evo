package info

import (
	"Evo/config"
	"Evo/service/game"
	"Evo/util"
	"github.com/gin-gonic/gin"
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
		RemainTime      float64
	}
	util.Success(c, "success", gin.H{
		"time": timeInfo{
			BeginTime:       config.StartTime.Unix(),
			EndTime:         config.EndTime.Unix(),
			NowRound:        config.ROUND_NOW,
			RoundTime:       config.ROUND_TIME,
			RoundRemainTime: config.GetRoundRemainTime(),
			RemainTime:      config.GetRestTime(),
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

// GetRank 获取排名
func GetRank(c *gin.Context) {
	util.Success(c, "success", gin.H{
		"rank": game.GetRankList(),
	})
}
