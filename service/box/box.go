/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		:对靶机状态的修改。
 */

package box

import (
	"Evo/db"
	"Evo/model"
)

// ResetAllStatus 重置靶机全部状态
func ResetAllStatus() error {
	err := db.DB.Model(&model.Box{}).Updates(map[string]interface{}{"is_attacked": false, "is_down": false}).Error
	if err != nil {
		return err
	}
	return nil
}

// ResetAllScore 重置靶机分数
func ResetAllScore() error {
	err := db.DB.Model(&model.Box{}).Update("score", db.DB.Model(&model.Challenge{}).
		Select("score").Where("boxes.challenge_id = challenges.id")).Error
	if err != nil {
		return err
	}
	return nil
}

func ResetAllAttack() error {
	if err := db.DB.Model(&model.Box{}).Where("is_down = ?", false).
		Update("is_attacked", false).Error; err != nil {
		return err
	}
	return nil
}