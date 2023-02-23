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
)

// CalcAttack 统计分数
func CalcAttack() error {
	round := config.ROUND_NOW - 1
	var allBoxes []model.GameBox
	// 直接查出所有靶机
	db.DB.Select([]string{"id", "challenge_id", "team_id", "score"}).Find(&allBoxes)

	var attacks []model.Attack
	// 把这一轮所有attack记录查出来
	db.DB.Where("round = ?", round).Find(&attacks)

	boxMap := make(map[uint]int)
	for i := 0; i < len(allBoxes); i++ {
		boxMap[allBoxes[i].ID] = i
	}
	// 被攻击的 靶机号：攻击他的靶机的靶机号
	attackMap := make(map[uint][]uint)
	for i := 0; i < len(attacks); i++ {
		attackMap[attacks[i].BoxId] = append(attackMap[attacks[i].BoxId], attacks[i].Attacker)
	}
	// 遍历被attack的靶机，更新靶机们的分数
	attackScore := float64(config.ATTACK_SCORE)
	for attacked, attackers := range attackMap {
		num := float64(len(attackers))
		allBoxes[boxMap[attacked]].Score -= attackScore
		for _, attacker := range attackers {
			allBoxes[boxMap[attacker]].Score += attackScore / num
		}
	}

	// 将靶机信息重新保存回去
	err := db.DB.Save(&allBoxes).Error
	return err
}

func CalcDown() error {
	downScore := float64(config.DOWN_SCORE)
	boxMap1 := make(map[uint][]uint)
	var downBoxes []model.GameBox
	db.DB.Where("is_down = ?", true).Select([]string{"id", "challenge_id", "team_id", "score"}).Find(&downBoxes)
	for i := 0; i < len(downBoxes); i++ {
		boxMap1[downBoxes[i].ChallengeID] = append(boxMap1[downBoxes[i].ChallengeID], downBoxes[i].ID)
	}

	boxMap2 := make(map[uint][]uint)
	var normalBoxes []model.GameBox
	db.DB.Where("is_down = ?", false).Select([]string{"id", "challenge_id", "team_id", "score"}).Find(&normalBoxes)
	for i := 0; i < len(normalBoxes); i++ {
		boxMap2[normalBoxes[i].ChallengeID] = append(boxMap2[normalBoxes[i].ChallengeID], normalBoxes[i].ID)
	}
	for challenge, boxes := range boxMap1 {
		for _, box := range boxes {
			downBoxes[box].Score -= downScore
		}
		num := float64(len(boxMap2[challenge]))
		for _, box := range boxMap2[challenge] {
			normalBoxes[box].Score += downScore / num
		}
	}

	err := db.DB.Save(&downBoxes).Error
	if err != nil {
		return err
	}
	err = db.DB.Save(&normalBoxes).Error
	if err != nil {
		return err
	}
	return nil
}
