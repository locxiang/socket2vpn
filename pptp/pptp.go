package pptp

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os/exec"
	"regexp"
	"socket2vpn/config"
	"socket2vpn/env"
	"socket2vpn/util"
	"strings"
	"time"
)

// 关闭所有管道
func CloseAll() error {
	log.Debug("关闭通道所有通道")
	cmd := exec.Command("poff", "-a")
	return cmd.Run()
}

func ClosePPTP(pptpName string) error {
	log.Debug("关闭通道：", pptpName)
	cmd := exec.Command("poff", pptpName)
	return cmd.Run()
}

// 建立pptp通道
func NewPPTP(u config.User) error {

	ClosePPTP(u.User)

	fmt.Printf("start pptp [%s] \n", u.User)

	serverIP := config.Values.Servers[rand.Intn(len(config.Values.Servers)-1)]
	cmdStr := fmt.Sprintf(" --creat %s --server %s --username %s --password %s --encrypt --start", u.User, serverIP, u.User, u.Pass)
	log.Debug("创建PPTP：", cmdStr)
	cmd := exec.Command("pptpsetup", strings.Split(cmdStr, " ")...)
	//cmd := exec.Command("ls", "-al")

	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Start()
	done := make(chan error, 1)

	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-time.After(10 * time.Second):
		cmd.Process.Kill()
	case err := <-done:
		if err != nil {
			log.Warn("执行出错: %v", err, errbuf.String())
		}
	}

	connectName, err := GetConnectName(outbuf.String())
	if err != nil {
		return err
	}

	//写入缓存
	log.Debug("写入pptp到缓存")
	ppp, err := getPPTP(u.User, connectName)
	if err != nil {
		return err
	}
	env.SavePPTP(u.User, ppp)

	return nil

}

func getPPTP(user, connectName string) (*env.PPTP, error) {
	ip, err := util.GetPPPIp(connectName)
	if err != nil {
		return nil, err
	}

	p := &env.PPTP{
		Ip:          ip,
		PPTPName:    user,
		ConnectName: connectName,
	}

	return p, nil

}

// 获取通道名称
func GetConnectName(str string) (string, error) {

	re, _ := regexp.Compile("ppp[0-9a-zA-Z_]+")
	one := re.Find([]byte(str))
	if len(one) == 0 {
		return "", errors.New("未找到通道")
	}
	return string(one), nil

}
