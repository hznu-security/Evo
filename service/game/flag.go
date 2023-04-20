/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		:比赛中flag相关的函数
 */

package game

import (
	"Evo/config"
	"Evo/db"
	"Evo/model"
	"Evo/util"
	"log"
	"strings"
)

func GenerateFlag() error {
	var boxes []model.GameBox
	db.DB.Find(&boxes)
	// 直接遍历靶机
	for _, box := range boxes {
		flags := make([]model.Flag, config.GAME_ROUND)
		for i := 0; i < int(config.GAME_ROUND); i++ {
			flags[i].GameBoxId = box.ID
			flags[i].ChallengeID = box.ChallengeID
			flags[i].Round = uint(i + 1)
			flags[i].TeamId = box.TeamId
			flags[i].Flag = util.GetRandomStr(15, config.FLAG_PRE, config.FLAG_SUF)
		}
		err := db.DB.Create(&flags).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func RefreshFlag(round uint) {
	var challenges []model.Challenge
	db.DB.Where("auto_refresh = ?", true).Find(&challenges)

	var flags []model.Flag
	db.DB.Where("round = ?", round).Find(&flags)
	var boxes []model.GameBox
	db.DB.Find(&boxes)

	m := make(map[uint]string)
	for _, flag := range flags {
		m[flag.GameBoxId] = flag.Flag
	}
	for _, challenge := range challenges {
		var gameBoxes []model.GameBox
		db.DB.Where("challenge_id = ?", challenge.ID).Find(&gameBoxes)
		for _, gameBox := range gameBoxes {
			go execRefresh(gameBox, m[gameBox.ID], challenge.Command, round)
		}
	}
}

func execRefresh(gameBox model.GameBox, flag string, command string, round uint) {
	command = strings.ReplaceAll(command, "{{FLAG}}", flag)
	_, err := util.SSHExec(gameBox.Port, gameBox.SshUser, gameBox.SshPwd, command)
	if err != nil {
		log.Printf("ssh error. Team:%d,GameBox:%s,Round:%d 更新FLAG失败", gameBox.TeamId, gameBox.ID, round)
	}
}
