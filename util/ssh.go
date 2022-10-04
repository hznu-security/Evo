/* COPYRIGHT NOTICE
 * 作者		:ymk
 * 创建时间	:2022/07/09 14:14
 * 描述		:ssh相关的工具
 */

package util

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"log"
	"time"
)

func SSHExec(ip string, user string, pwd string, cmd string) (output string, err error) {
	client, err := ssh.Dial("tcp", ip+":22", &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pwd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	})
	if err != nil {
		return "", err
	}

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	var stdout bytes.Buffer
	session.Stdout = &stdout

	err = session.Run(cmd)
	if err != nil {
		return "", err
	}
	err = client.Close()
	if err != nil {
		log.Println(err)
	}

	return stdout.String(), nil
}
