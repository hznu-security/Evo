/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		:测试util包里的函数
 */

package util

import (
	"log"
	"testing"
)

// 测试SSHExec
func TestSSHExec(t *testing.T) {
	ip := "101.37.13.45"
	user := "ymk"
	pwd := "ymk0910"
	cmd := "whoami"
	res, err := SSHExec(ip, user, pwd, cmd)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res)
}
