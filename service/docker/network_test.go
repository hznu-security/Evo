/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/23 11:14
 * 描述     ：测试网络相关函数
 */

package docker

import "testing"

// 测试创建一个名为test的网络
func TestCreateNetwork(t *testing.T) {
	name := "test"
	subnet := "172.18.0.0/16"
	res, err := CreateNetwork(name, subnet)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	t.Log(res)
}

//测试删除一个名为test的网络
func TestDelNetwork(t *testing.T) {
	name := "test"
	err := DelNetwork(name)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
}
