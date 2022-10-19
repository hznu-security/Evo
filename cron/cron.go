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
	"github.com/robfig/cron/v3"
	"log"
	"strconv"
	"strings"
	"time"
)

var startEntry, endEntry, roundEntry cron.EntryID // 一个开启比赛的计时任务,一个结束比赛的计时任务,一个每轮比赛执行一次的计时任务
var cron1 = cron.New()

// StartCron 开启计时
func StartCron() error {
	specStart, err := parseTime(config.START_TIME)
	if err != nil {
		return errors.New("解析 start_time 失败")
	}
	specEnd, err := parseTime(config.END_TIME)
	if err != nil {
		return errors.New("解析 start_time 失败")
	}

	specRound := "@every" + strconv.Itoa(int(config.ROUND_TIME)) + "m"

	if err != nil {
		panic(err)
	}

	startEntry, _ = cron1.AddFunc(specStart, func() {
		roundEntry, _ = cron1.AddFunc(specRound, func() { // 注册每10分钟一执行的任务
			log.Println("新回合开始")
			config.ROUND_NOW++
			game.RefreshFlag()
		})
	})
	endEntry, _ = cron1.AddFunc(specEnd, func() {
		config.ROUND_NOW++ // 确保时间到了之后，ROUND_NOW的值一定大于GAME_ROUND
		cron1.Remove(startEntry)
		cron1.Remove(roundEntry)
		cron1.Remove(endEntry)
	})
	cron1.Start()
	return nil
}

// TerminateCron 终止计时
func TerminateCron() {
	cron1.Remove(startEntry)
	cron1.Remove(endEntry)
	cron1.Remove(roundEntry)
	cron1.Stop()
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
