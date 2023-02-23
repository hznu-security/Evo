/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		:缓存
 */

package db

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var store *cache.Cache

func Init() {
	c := cache.New(cache.NoExpiration, cache.DefaultExpiration)
	store = c
}

func Get(k string) (interface{}, bool) {
	return store.Get(k)
}

func Set(k string, x interface{}, d ...time.Duration) {
	duration := cache.NoExpiration
	if len(d) == 1 {
		duration = d[0]
	}
	store.Set(k, x, duration)
}
