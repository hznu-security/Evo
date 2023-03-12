/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		:处理文件路径
 */

package util

import (
	"log"
	"os"
)

// TestAndMake 检验路径是否存在，不存在就创建它
func TestAndMake(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		log.Printf("路径不村子啊:%s", path)
		err := os.MkdirAll(path, 0775)
		if err != nil {
			return err
		}
		log.Printf("创建路径:%s", path)
	}
	return nil
}
