/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		:管理比赛配置
 */

package ctrl

import (
	"Evo/cron"
	"Evo/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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
		str := v.(string)
		if v == "" || len(str) > 50 {
			util.Fail(c, "参数:"+k+"异常", nil)
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
			viper.Set(key, v)
		}
	}

	// 覆盖之前的配置
	err := viper.WriteConfig()
	if err != nil {
		util.Error(c, "配置写入失败", nil)
		return
	}
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
