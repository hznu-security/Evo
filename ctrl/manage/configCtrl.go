/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		:管理比赛配置
 */

package manage

import (
	"Evo/config"
	"Evo/cron"
	"Evo/db"
	"Evo/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"time"
)

// PutConfig 修改比赛相关配置
func PutConfig(c *gin.Context) {
	// value必须是interface{}类型，不然bool类型绑定不上
	var m = make(map[string]interface{})
	if err := json.NewDecoder(c.Request.Body).Decode(&m); err != nil {
		util.Fail(c, "参数绑定异常", gin.H{
			"config": m,
		})
		return
	}

	if len(m) > 20 {
		util.Fail(c, "参数项过多，请检查", nil)
		return
	}

	for k, v := range m {
		_, ok := v.(bool)
		if ok {
			continue
		}
		str, ok := v.(string)
		if ok {
			if str == "" || len(str) > 50 {
				util.Fail(c, fmt.Sprintf("参数 %s 异常 %v", k, v), nil)
				return
			} else {
				continue
			}
		}
		_, ok = v.(int)
		if ok {
			continue
		} else {
			util.Fail(c, fmt.Sprintf("参数 %s 异常 %v", k, v), nil)
			return
		}

	}

	for k, v := range m {
		key := "game." + k
		val := viper.GetString(key)
		if val == "" {
			util.Fail(c, "配置项:"+k+"不存在", nil)
			return
		} else {
			if value, ok := v.(int); ok {
				viper.Set(key, value)
				continue
			}
			if value, ok := v.(bool); ok {
				viper.Set(key, value)
				continue
			}
			if value, ok := v.(string); ok {
				viper.Set(key, value)
				continue
			}
		}
	}
	// 覆盖之前的配置
	err := viper.WriteConfig()
	if err != nil {
		util.Error(c, "配置写入失败", nil)
		return
	}

	st := viper.GetString("game.starttime")
	start, err := time.ParseInLocation(config.TIME_FORMAT, st, time.Local)
	if time.Now().Sub(start) > 0 {
		util.Success(c, "比赛开始时间早于当前时间", nil)
	}
	//config.SetTime()
	util.Success(c, "success", nil)
}

// GetConfig 获取比赛相关配置
func GetConfig(c *gin.Context) {
	configs := viper.GetStringMapString("game")
	util.Success(c, "success", gin.H{
		"configs": configs,
	})
}

// ResetConfig 重置平台  感觉没啥好写的
func ResetConfig(c *gin.Context) {
	// 删除team表,box表,challenges表，flags表，game_box表，notifications表，webhooks表,attack表，down表，chat表
	sql := fmt.Sprintf("TRUNCATE TABLE %s;", "downs")
	log.Println(sql)
	db.DB.Exec(sql)

	sql = fmt.Sprintf("TRUNCATE TABLE %s;", "attacks")
	log.Println(sql)
	db.DB.Exec(sql)

	sql = fmt.Sprintf("TRUNCATE TABLE %s;", "teams")
	log.Println(sql)
	db.DB.Exec(sql)

	sql = fmt.Sprintf("truncate table %s;", "boxes")
	log.Println(sql)
	db.DB.Exec(sql)

	sql = fmt.Sprintf("TRUNCATE TABLE %s;", "challenges")
	log.Println(sql)
	db.DB.Exec(sql)

	sql = fmt.Sprintf("TRUNCATE TABLE %s;", "game_boxes")
	log.Println(sql)
	db.DB.Exec(sql)

	sql = fmt.Sprintf("TRUNCATE TABLE %s;", "flags")
	log.Println(sql)
	db.DB.Exec(sql)

	sql = fmt.Sprintf("TRUNCATE TABLE %s;", "notifications")
	log.Println(sql)
	db.DB.Exec(sql)

	sql = fmt.Sprintf("TRUNCATE TABLE %s;", "webhooks")
	log.Println(sql)
	db.DB.Exec(sql)

	sql = fmt.Sprintf("TRUNCATE TABLE %s;", "charts")
	log.Println(sql)
	db.DB.Exec(sql)
}

// StartGame 设置比赛可以开始，到时间自动开启比赛
func StartGame(c *gin.Context) {
	err := cron.StartCron()
	if err != nil {
		util.Error(c, err.Error(), nil)
		return
	}
	util.Success(c, "success", nil)
}

// TerminateGame 终止比赛
func TerminateGame(c *gin.Context) {
	cron.TerminateCron()
	util.Success(c, "success", nil)
}
