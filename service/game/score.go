/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		:分数结算相关
 */

package game

import (
	"Evo/config"
	"Evo/db"
	"Evo/model"
	"log"
	"time"
)

// CalcScore 统计该回合的分数
func CalcScore(round uint) {
	if round <= 0 {
		return
	}
	start := time.Now()
	addAttack(round)
	minusAttack(round)
	minusCheckDown(round)
	addCheckDown(round)
	calcGameBoxScore()
	calcTeamScore()
	end := time.Now()
	log.Printf("分数结算完成,耗时%f\n", end.Sub(start).Seconds())
}

func addAttack(round uint) {
	var gameBoxes []model.GameBox
	db.DB.Find(&gameBoxes)
	for _, gameBox := range gameBoxes {
		var attacks []model.Attack
		db.DB.Where(&model.Attack{GameBoxId: gameBox.ID, Round: round}).Find(&attacks)
		if len(attacks) != 0 {
			score := float64(config.ATTACK_SCORE) / float64(len(attacks))
			scores := make([]model.Score, 0)
			for _, attack := range attacks {
				var attackerGameBox model.GameBox
				db.DB.Where(&model.GameBox{ChallengeID: attack.ChallengeId, TeamId: attack.Attacker}).First(&attackerGameBox)
				scores = append(scores, model.Score{
					TeamId:    attack.Attacker,
					GameBoxId: attackerGameBox.ID,
					Round:     round,
					Reason:    "attack",
					Score:     score,
				})
			}
			db.DB.Create(&scores)
		}
	}
}

func minusAttack(round uint) {
	var attacks []struct {
		GameBoxId uint `gorm:"game_box_id"`
		TeamId    uint `gorm:"team_id"`
	}
	db.DB.Table("attacks").Select("DISTINCT(`game_box_id`) AS game_box_id,team_id").
		Where("round = ?", round).Scan(&attacks)
	scores := make([]model.Score, 0)
	for _, attack := range attacks {
		scores = append(scores, model.Score{
			TeamId:    attack.TeamId,
			GameBoxId: attack.GameBoxId,
			Round:     round,
			Reason:    "been_attacked",
			Score:     float64(-config.ATTACK_SCORE),
		})
	}
	db.DB.Create(&scores)
}

func minusCheckDown(round uint) {
	var downs []model.Down
	db.DB.Where("round = ?", round).Find(&downs)

	var scores []model.Score
	for _, down := range downs {
		scores = append(scores, model.Score{
			TeamId:    down.TeamId,
			GameBoxId: down.GameBoxId,
			Round:     round,
			Reason:    "checkdown",
			Score:     float64(-config.DOWN_SCORE),
		})
	}
	db.DB.Create(&scores)
}

func addCheckDown(round uint) {
	var challenges []model.Challenge
	db.DB.Find(&challenges)
	for _, challenge := range challenges {
		var downs []model.Down
		db.DB.Where("challenge_id = ? AND round = ?", challenge.ID, round).Find(&downs)
		// 总的扣分
		totalScore := len(downs) * int(config.DOWN_SCORE)

		// 被扣分的靶机
		var downGameBox []uint
		for _, down := range downs {
			downGameBox = append(downGameBox, down.GameBoxId)
		}

		// 查出加分的靶机
		var addGameBox []model.GameBox
		db.DB.Where("challenge_id = ?", challenge.ID).Not("id", downGameBox).Find(&addGameBox)

		var scores []model.Score
		score := float64(totalScore) / float64(len(addGameBox))
		for _, gameBox := range addGameBox {
			scores = append(scores, model.Score{
				TeamId:    gameBox.TeamId,
				GameBoxId: gameBox.ID,
				Round:     round,
				Score:     score,
				Reason:    "service_online",
			})
		}
		db.DB.Create(scores)
	}
}

func calcGameBoxScore() {
	var gameBoxes []model.GameBox
	db.DB.Find(&gameBoxes)
	for _, gameBox := range gameBoxes {
		var score struct {
			Score float64 `gorm:"Column:Score"`
		}
		db.DB.Table("scores").Select("SUM(score) AS Score").
			Where("game_box_id = ?", gameBox.ID).Scan(&score)
		var challenge model.Challenge
		db.DB.Where("id = ?", gameBox.ChallengeID).First(&challenge)                   // 获得基础分数
		db.DB.Where("id = ?", gameBox.ID).Update("score", challenge.Score+score.Score) //计算靶机分数
	}
}

func calcTeamScore() {
	var teams []model.Team
	db.DB.Find(&teams)
	for _, team := range teams {
		var score struct {
			Score float64 `gorm:"Score"`
		}
		db.DB.Select("SUM(score) AS Score").Where("team_id = ? AND visible = ?", team.ID, true).Scan(&score)
		db.DB.Where("id = ?", team.ID).Update("score", score.Score)
	}
}
