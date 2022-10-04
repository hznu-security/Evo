/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		:从配置文件初始化一些配置
 */

package config

import (
	"github.com/spf13/viper"
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

//InitConfig 初始化配置
func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
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
	start, err := time.ParseInLocation(TIME_FORMAT, START_TIME, time.Local)
	if err != nil {
		panic("加载比赛时间失败")
	}
	end, err := time.ParseInLocation(TIME_FORMAT, END_TIME, time.Local)
	if err != nil {
		panic("加载比赛时间失败")
	}
	processing := uint(start.Sub(end).Minutes())
	GAME_ROUND = processing / ROUND_TIME
}
