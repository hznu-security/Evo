/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		:实现平台中计时的模块
 */

package cron

import (
	"Evo/config"
	"Evo/service/game"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

func InitCron() {
	cron1 := cron.New()

	// var entryEnd cron.EntryID
	var entryStart, entryRound cron.EntryID
	specStart, err := parseTime(config.START_TIME)
	if err != nil {
		panic(errors.New("解析 start_time 失败"))
	}
	specEnd, err := parseTime(config.END_TIME)
	if err != nil {
		panic(errors.New("解析 start_time 失败"))
	}

	specRound := "@every" + strconv.Itoa(int(config.ROUND_TIME)) + "m"

	if err != nil {
		panic(err)
	}
	entryStart, _ = cron1.AddFunc(specStart, func() {
		_, _ = cron1.AddFunc(specRound, func() { // 注册每10分钟一执行的任务
			log.Println("新回合开始")
			config.ROUND_NOW++
			game.RefreshFlag()
		})
	})
	_, _ = cron1.AddFunc(specEnd, func() {
		cron1.Remove(entryRound)
		cron1.Remove(entryStart)
	})

}

func parseTime(timeStr string) (string, error) {
	t, err := time.ParseInLocation(config.TIME_FORMAT, timeStr, time.Local)
	if err != nil {
		return "", err
	}
	builder := strings.Builder{}
	min := t.Minute()
	hour := t.Hour()
	year, month, day := t.Date()
	builder.WriteString(strconv.Itoa(min))
	builder.WriteString(" ")
	builder.WriteString(strconv.Itoa(hour))
	builder.WriteString(" ")
	builder.WriteString(strconv.Itoa(day))
	builder.WriteString(" ")
	builder.WriteString(strconv.Itoa(int(month)))
	builder.WriteString(" ")
	builder.WriteString(strconv.Itoa(year))
	return builder.String(), nil
}
