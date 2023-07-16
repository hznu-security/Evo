/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		:测试flag模块
 */

package game

import (
	"Evo/config"
	"Evo/db"
	"testing"
)

func BenchmarkFlag(b *testing.B) {
	config.InitConfig()
	db.InitDB()
	// 每次都刷第一回合的进去
	for n := 0; n < b.N; n++ {
		RefreshFlag(1)
	}
}
