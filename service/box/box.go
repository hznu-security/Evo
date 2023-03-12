/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		:对靶机状态的修改。
 */

package box

import (
	"Evo/db"
	"Evo/model"
	"errors"
	"gorm.io/gorm"
)

// ResetAllStatus 重置靶机全部状态
func ResetAllStatus() error {
	err := db.DB.Model(&model.GameBox{}).Updates(map[string]interface{}{"is_attacked": false, "is_down": false}).Error
	if err != nil && !errors.Is(err, gorm.ErrMissingWhereClause) {
		return err
	}
	return nil
}

// ResetAllScore 重置靶机分数  // TODO
func ResetAllScore() error {
	var challenges []model.Challenge
	db.DB.Select([]string{"id", "score"}).Find(&challenges)
	challengeMap := make(map[uint]float64)
	for _, challenge := range challenges {
		challengeMap[challenge.ID] = challenge.Score
	}
	var boxes []model.GameBox
	db.DB.Model(model.GameBox{}).Where("is_attacked = ? OR is_down = ?", true, true).Find(&boxes)
	for i := 0; i < len(boxes); i++ {
		boxes[i].Score = challengeMap[boxes[i].ChallengeID]
		boxes[i].IsDown = false
		boxes[i].IsAttacked = false
	}
	err := db.DB.Save(&boxes).Error
	if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
		return err
	}
	return nil
}
