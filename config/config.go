/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		:从配置文件初始化一些配置
 */

package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

const TIME_FORMAT = "2006-01-02 15:04:05"

//比赛相关设置

// GAME_NAME 比赛名字
var GAME_NAME string

// ROUND_TIME 每轮的时间，分钟为单位
var ROUND_TIME uint

// GAME_ROUND 比赛轮次
var GAME_ROUND uint

// FLAG_PRE flag前缀
var FLAG_PRE string

// FLAG_SUF flag 后缀
var FLAG_SUF string

// STARRY_ON 大屏是否开启
var STARRY_ON bool

// BOX_VISIBLE 队伍靶机互相可见
var BOX_VISIBLE bool

// START_TIME 比赛开始时间
var START_TIME string

// END_TIME 比赛结束时间
var END_TIME string

// ROUND_NOW 当前比赛的轮次
var ROUND_NOW uint

// DOWN_SCORE checkdown 扣分
var DOWN_SCORE uint

// ATTACK_SCORE 被攻击扣分
var ATTACK_SCORE uint

var StartTime time.Time

var EndTime time.Time

func SetTime() {
	start, err := time.ParseInLocation(TIME_FORMAT, START_TIME, time.Local)
	if err != nil {
		panic("加载比赛时间失败,检查时间格式")
	}
	end, err := time.ParseInLocation(TIME_FORMAT, END_TIME, time.Local)
	if err != nil {
		panic("加载比赛时间失败,检查时间格式")
	}

	if end.Sub(start) < 0 {
		log.Panicln("比赛时间设置有误,结束时间应晚于开始时间")
	}

	StartTime = start
	EndTime = end
	processing := uint(end.Sub(start).Minutes()) // 后面的时间 sub 前面的时间
	GAME_ROUND = processing / ROUND_TIME
}

func SetTime1() {
	start, err := time.ParseInLocation(TIME_FORMAT, START_TIME, time.Local)
	if err != nil {
		panic("加载比赛时间失败,检查时间格式")
	}
	end, err := time.ParseInLocation(TIME_FORMAT, END_TIME, time.Local)
	if err != nil {
		panic("加载比赛时间失败,检查时间格式")
	}

	StartTime = start
	EndTime = end
	if EndTime.Sub(StartTime) < 0 {
		log.Panicln("开始时间应早于结束时间")
	}
	processing := uint(end.Sub(start).Minutes()) // 后面的时间 sub 前面的时间
	if processing%ROUND_TIME != 0 {
		log.Panicln("回合数不为整数，请重设时间")
	}
	GAME_ROUND = processing / ROUND_TIME

	// 获得round now
	ROUND_NOW = uint(time.Now().Sub(StartTime) / time.Minute * time.Duration(ROUND_TIME))
}

// InitConfig 初始化配置
func InitConfig() {
	log.Println("读取设置")
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	GAME_NAME = viper.GetString("game.name")
	FLAG_PRE = viper.GetString("game.prefix")
	FLAG_SUF = viper.GetString("game.suffix")
	STARRY_ON = viper.GetBool("game.starry")
	BOX_VISIBLE = viper.GetBool("game.boxvisible")
	START_TIME = viper.GetString("game.starttime")
	END_TIME = viper.GetString("game.endtime")
	DOWN_SCORE = viper.GetUint("game.downscore")
	ATTACK_SCORE = viper.GetUint("game.attackscore")
	ROUND_TIME = viper.GetUint("game.roundtime")

	//SetTime()

	SetTime1()
}

// 返回本轮剩余时间 返回秒数
func GetRoundRemainTime() float64 {
	res := StartTime.Add(time.Duration(ROUND_NOW*ROUND_TIME) * time.Minute).Sub(time.Now()).Seconds()
	if res > 0 {
		return res
	} else {
		return 0
	}
}

// 返回比赛剩余时间  返回秒数
func GetRestTime() float64 {
	res := EndTime.Sub(time.Now()).Seconds()
	if res > 0 {
		return res
	} else {
		return 0
	}
}
