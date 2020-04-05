package env

import (
	"net"
	"sync"
)

var list sync.Map

type PPTP struct {
	Port        int //socket5 端口号
	Ip          *net.TCPAddr
	PPTPName    string //拨号名称  就是用户名
	ConnectName string //管道名称  ppp0
}

// 添加
func SavePPTP(user string, c *PPTP) {
	list.Store(user, c)
}

// 更新端口
func SetSocket5Port(user string, port int) {
	v, _ := list.Load(user)
	pp := v.(*PPTP)

	pp.Port = port
}

func GetPPTP(user string) *PPTP {
	v, _ := list.Load(user)
	return v.(*PPTP)
}
