/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/10 13:15
 * 描述     ：测试jwt模块
 */

package auth

import (
	"fmt"
	"testing"
	"time"
)

func TestReleaseAndParse(t *testing.T) {
	var id uint = 10
	role := TEAM
	tokenString, err := ReleaseToken(id, role)
	t.Log("token:", tokenString)
	if err != nil {
		fmt.Println(err.Error())
	}
	time.Sleep(time.Second * 5)
	token, _, err := ParseToken(tokenString)
	if err != nil {
		t.Log(err.Error())
	}
	if !token.Valid {
		t.Log("错了")
	}
	tokenString = tokenString + "a"
	token, _, err = ParseToken(tokenString)
	if err != nil {
		t.Log(err.Error())
	}
	if token.Valid {
		t.Log("错了")
	}
}
