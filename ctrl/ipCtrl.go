/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		:管理主机上的ip地址
 */

package ctrl

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"os/exec"
)

/**
已弃用
*/

func GetInterfaces(c *gin.Context) {
	ifis, err := net.Interfaces()
	if err != nil {
		log.Printf("获取网卡失败. %v\n", err)
		Error(c, fmt.Sprintf("获取网卡失败. %s", err.Error()), nil)
		return
	}
	ifiNames := make([]string, len(ifis))
	for i := 0; i < len(ifis); i++ {
		ifiNames[i] = ifis[i].Name
	}
	Success(c, "success", gin.H{
		"ifis": ifiNames,
	})
}

// GetIpAddress 默认拿ens33的ip
func GetIpAddress(c *gin.Context) {
	interfaceName := c.Query("ifi")
	if interfaceName == "" {
		Fail(c, BindError, nil)
		return
	}
	addresses, err := getIpAddress(interfaceName)
	if err != nil {
		log.Println(err)
		Error(c, err.Error(), nil)
		return
	} else {
		Success(c, "success", gin.H{
			"ipList": addresses,
		})
		return
	}
}

type ipForm struct {
	Ifi string `binding:"required,max=20"`
	Ip  string `binding:"required,max=100"`
}

// PostIpAddress 添加Ip地址
func PostIpAddress(c *gin.Context) {
	var form ipForm
	if err := c.ShouldBind(&form); err != nil {
		Fail(c, BindError, nil)
		return
	}
	fmt.Println(form)
	_, err := addIpAddress(form.Ifi, form.Ip)
	if err != nil {
		log.Println(err.Error())
		Error(c, fmt.Sprintf("添加ip地址出错. %v", err), nil)
		return
	}
	Success(c, "success", nil)
	return
}

func DelIpAddress(c *gin.Context) {
	var form ipForm
	if err := c.ShouldBind(&form); err != nil {
		Fail(c, BindError, nil)
	}
	_, err := delIpAddress(form.Ifi, form.Ip)
	if err != nil {
		log.Println(err)
		Error(c, fmt.Sprintf("删除ip地址出错. %v", err), nil)
		return
	}
	Success(c, "success", nil)
}

func getIpAddress(interfaceName string) ([]string, error) {
	ifi, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return nil, fmt.Errorf("打开网卡错误. %v", err)
	}
	addrs, err := ifi.Addrs()
	if err != nil {
		return nil, fmt.Errorf("查找网卡ip出错. %v", err)
	}
	res := make([]string, 0)
	for _, addr := range addrs {
		if ipv4Addr := addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
			res = append(res, ipv4Addr.String())
		}
	}
	return res, nil
}

func addIpAddress(interfaceName string, ipAddress string) (string, error) {
	cmd := exec.Command("ip", "address", "add", ipAddress, "dev", interfaceName)
	output, err := cmd.Output()
	if err != nil {
		return string(output), fmt.Errorf("添加ip地址出错. %s", string(err.(*exec.ExitError).Stderr))
	}
	return string(output), nil
}

func delIpAddress(interfaceName string, ipAddress string) (string, error) {
	cmd := exec.Command("ip", "address", "del", ipAddress, "dev", interfaceName)
	output, err := cmd.Output()
	if err != nil {
		return string(output), fmt.Errorf("删除ip地址出错. %s", string(err.(*exec.ExitError).Stderr))
	}
	return string(output), nil
}
