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
	"math/rand"
	"strings"
	"time"
)

// GenerateFlag 生成flag
func GenerateFlag() error {
	var boxes []model.Box
	db.DB.Find(&boxes)
	// 直接遍历靶机
	for _, box := range boxes {
		flags := make([]model.Flag, config.GAME_ROUND)
		for i := 0; i < int(config.GAME_ROUND); i++ {
			flags[i].BoxId = box.ID
			flags[i].ChallengeID = box.ChallengeID
			flags[i].Round = uint(i + 1)
			flags[i].TeamId = box.TeamId
			flags[i].Flag = generateFlag(config.FLAG_PRE, config.FLAG_SUF)
		}
		err := db.DB.Create(&flags).Error
		if err != nil {
			return err
		}
	}
	return nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano()) // 使用自己的source，不考虑并发安全，性能提升

func generateFlag(prefix, suffix string) string {
	n := 15 + len(prefix) + len(suffix)
	sb := strings.Builder{} // stringbuilder，提高拼接字符串的效率
	sb.Grow(n)
	sb.WriteString(prefix)
	// 掩码，多次利用，提高性能。
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	sb.WriteString(suffix)
	return sb.String()
}

func RefreshFlag() {
	/**
	为每一台靶机更新flag
	更新过程：
		查出flag      boxId:flag
		查出靶机信息   boxId:boxInfo
		执行ssh
	*/
	var flags []model.Flag
	db.DB.Where("round = ?", config.ROUND_NOW).Find(&flags)
	var boxes []model.Box
	db.DB.Find(&boxes)

	m := make(map[uint]string)
	for _, flag := range flags {
		m[flag.BoxId] = flag.Flag
	}
	for i := 0; i < len(boxes); i++ {
		go execRefresh(&boxes[i], m[boxes[i].ID])
	}

}

func execRefresh(box *model.Box, flag string) {

}
