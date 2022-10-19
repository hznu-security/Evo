/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/10 12:53
 * 描述     ：这里是处理密码的函数，对密码进行编码，比较密码是否正确,为队伍生成随机的密码
 */

package auth

import (
	"Evo/util"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// Encode 对密码进行编码
func Encode(pwd string) string {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	return string(hashedPwd)
}

// Cmp 比较密码是否正确
func Cmp(epwd string, pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(epwd), []byte(pwd)) == nil
}

// NewPwd 生成队伍密码
func NewPwd() string {
	return util.GetRandomStr(8, "", "")
}
