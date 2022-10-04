/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/23 11:15
 * 描述     ：测试容器相关方法
 */
package docker

import (
	"bytes"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

func TestStartContainer(t *testing.T) {
	image := "easyweb"
	ip := "172.18.0.6"
	name := "testweb"
	net := "test"
	err := StartContainer(image, name, net, ip)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
}

// 设置testweb这个容器地密码，并且进去执行whoami
func TestSetContainerSSH(t *testing.T) {
	container := "testweb"
	user := "root"
	pwd := "123456"
	err := SetContainerSSH(container, user, pwd)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	client, err := ssh.Dial("tcp", "172.18.0.6:22", &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.Password(pwd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	})
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	var output bytes.Buffer
	session.Stdout = &output

	err = session.Run("whoami")
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	t.Log(output.String())
}

func TestRemoveContainer(t *testing.T) {
	name := "testweb"
	if err := RemoveContainer(name); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestResetContainer(t *testing.T) {
	name := "testweb"
	pwd, err := ResetContainer(name)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(pwd)
}

func TestContainerExec(t *testing.T) {
	cmd := "whoami"
	container:= "7ad5f1f1a149"
	res,err := ContainerExec(container,cmd)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(string(res))
}