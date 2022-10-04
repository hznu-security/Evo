/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/10 13:15
 * 描述     ：测试密码编码和比较是否正确
 */

package auth

import "testing"

func TestPwd(t *testing.T) {
	pwd := "123456"
	epwd := Encode(pwd)
	t.Log("encrypted pwd", epwd)
	wpwd := "1234567"
	res1 := Cmp(epwd, pwd)
	res2 := Cmp(epwd, wpwd)

	if !(res1 && !res2) {
		t.Fail()
	}
}
