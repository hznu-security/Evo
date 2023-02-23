/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/09 20:19
 * 描述     ：测试config能否初始化
 */

package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"testing"
	"time"
)

/**
执行单元测试需要在config目录下单独添加config.yml文件
*/

func TestConfig(t *testing.T) {
	//获取一个路径,到当前目录
	wd, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(wd)
	err := viper.ReadInConfig()
	if err != nil {
		t.Log(err.Error())
	}
	//读取config中关于数据库的设置
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	pwd := viper.GetString("datasource.pwd")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utfmb4&parseTime=True&loc=Local",
		username, pwd, host, port, database)
	t.Log(dsn)
}

func TestSetConfig(t *testing.T) {
	wd, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(wd)
	err := viper.ReadInConfig()
	if err != nil {
		t.Log(err.Error())
	}
	flag := viper.GetString("game.flag")
	score := viper.GetString("game.score")
	t.Log(flag, score)
	viper.Set("game.score", 10)
	viper.Set("game.flag", 123)
	flag = viper.GetString("game.flag")
	score = viper.GetString("game.score")
	t.Log(flag, score)
	err = viper.WriteConfig()
	if err != nil {
		t.Log(err)
	}
}

func TestGetMap(t *testing.T) {
	wd, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(wd)
	err := viper.ReadInConfig()
	if err != nil {
		t.Log(err.Error())
	}
	m := viper.GetStringMapString("game")
	fmt.Println(m)
}

func TestSetTime(t *testing.T) {
	wd, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(wd)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
	now := time.Now()
	viper.Set("game.starttime", now.Format("2006-01-02 15:04:05"))
	log.Println(now.Format("2006-01-02 15:04:05"))
	err = viper.WriteConfig()
	if err != nil {
		log.Fatalln(err)
	}
}

func TestTime(t *testing.T) {
	wd, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(wd)
	viper.ReadInConfig()
	startTime := viper.GetString("game.starttime")
	log.Println(startTime)
	start, err := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(start.Format("2006-01-02 15:04:05"))
}

// 测试计算比赛之间的轮数
func TestRound(t *testing.T) {
	wd, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(wd)
	viper.ReadInConfig()

	roundTime := viper.GetInt("game.roundtime") //每轮的时间
	log.Println(roundTime)
	startTime := viper.GetString("game.starttime")
	endTime := viper.GetString("game.endtime")
	start, _ := time.ParseInLocation(TIME_FORMAT, startTime, time.Local)
	end, _ := time.ParseInLocation(TIME_FORMAT, endTime, time.Local)
	processing := int(end.Sub(start).Minutes())
	rounds := processing / roundTime
	log.Println(rounds)
}

// 基准测试获取时间的操作，很快
func BenchmarkGetRoundRemainTime(b *testing.B) {
	InitConfig()
	for i := 0; i < b.N; i++ {
		GetRoundRemainTime()
	}
}
