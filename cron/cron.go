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
	// 获得开始计时的时间
	var cronStartTime time.Time
	var specStart string

	// 如果比赛时间晚于当前
	if time.Now().Sub(config.StartTime) < 0 {
		cronStartTime = config.StartTime
		config.ROUND_NOW = 0
	} else { // 比赛时间早于当前
		cronStartTime = config.StartTime.Add(time.Duration(config.ROUND_NOW * config.ROUND_TIME))
		config.ROUND_NOW = uint(time.Now().Sub(config.StartTime) / time.Minute * time.Duration(config.ROUND_TIME))
	}

	specStart, err := parseTime(cronStartTime)
	if err != nil {
		return errors.New("解析 start_time 失败")
	}

	specEnd, err := parseTime(config.EndTime)
	if err != nil {
		return errors.New("解析 start_time 失败")
	}

	specRound := "@every" + strconv.Itoa(int(config.ROUND_TIME)) + "m"

	if err != nil {
		panic(err.Error())
	}

	log.Println("初始化定时任务")
	startEntry, _ = cron1.AddFunc(specStart, func() {
		// 如果是比赛中途启动平台，先把flag刷进去
		if config.ROUND_NOW != 0 {
			game.RefreshFlag(config.ROUND_NOW)
			game.CalcScore(config.ROUND_NOW - 1)
		}
		roundEntry, _ = cron1.AddFunc(specRound, func() { // 注册每10分钟一执行的任务
			log.Println("新回合开始")
			config.ROUND_NOW++
			game.RefreshFlag(config.ROUND_NOW)
			game.CalcScore(config.ROUND_NOW - 1)
		})
	})
	endEntry, _ = cron1.AddFunc(specEnd, func() {
		config.ROUND_NOW++ // 确保时间到了之后，ROUND_NOW的值一定大于GAME_ROUND
		cron1.Remove(startEntry)
		cron1.Remove(roundEntry)
		cron1.Stop()
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

func parseTimeStr(timeStr string) (string, error) {
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

func parseTime(t time.Time) (string, error) {
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

func IsOn() bool {
	return config.ROUND_NOW <= config.GAME_ROUND && config.ROUND_NOW > 0
}
