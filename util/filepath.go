/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		:处理文件路径
 */

package util

import "os"

// TestAndMake 检验路径是否存在，不存在就创建它
func TestAndMake(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, 0775)
		if err != nil {
			return err
		}
	}
	return nil
}
