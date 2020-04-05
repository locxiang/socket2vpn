package proxy

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"socket2vpn/config"
	"socket2vpn/env"
	"strconv"
)

var (
	noAuth   = []byte{0x05, 0x00}
	withAuth = []byte{0x05, 0x02}

	authSuccess = []byte{0x05, 0x00}
	authFailed  = []byte{0x05, 0x01}

	connectSuccess = []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

type Socks5ProxyHandler struct {
	Port int //socket的端口
	Auth bool
	User string
	Pass string
}

type Handler interface {
	Handle(connect net.Conn)
}

func (socks5 Socks5ProxyHandler) getInfo(b []byte, connect net.Conn) (host, port string) {
	n, err := connect.Read(b)
	if err != nil {
		panic("getInfo失败：" + err.Error())
	}
	switch b[3] {
	case 0x01: //IP V4
		host = net.IPv4(b[4], b[5], b[6], b[7]).String()
	case 0x03: //domain
		host = string(b[5 : n-2]) //b[4] length of domain
	case 0x04: //IP V6
		host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
	default:
		return
	}
	port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))
	return host, port
}

func (socks5 *Socks5ProxyHandler) Handle(connect net.Conn) {
	if err := recover(); err != nil {
		log.Fatalf("err: %s", err)
		return
	}

	if connect == nil {
		return
	}
	defer connect.Close()

	b := make([]byte, 1024)

	_, err := connect.Read(b)
	if err != nil {
		return
	}

	if b[0] != 0x05 {
		return
	}

	if socks5.Auth == false {
		connect.Write(noAuth)
	} else {
		connect.Write(withAuth)

		_, err = connect.Read(b)
		if err != nil {
			return
		}

		userLength := int(b[1])
		user := string(b[2:(2 + userLength)])
		pass := string(b[(2 + userLength):])

		if socks5.User == user && socks5.Pass == pass {
			connect.Write(authSuccess)
		} else {
			connect.Write(authFailed)
			return
		}
	}

	pptp := env.GetPPTP(socks5.User)

	dialer := net.Dialer{
		LocalAddr: pptp.Ip,
	}

	host, port := socks5.getInfo(b, connect)
	server, err := dialer.Dial("tcp", net.JoinHostPort(host, port))
	if server != nil {
		defer server.Close()
	}
	if err != nil {
		return
	}
	connect.Write(connectSuccess)

	go io.Copy(server, connect)
	io.Copy(connect, server)
}

func NewSocket5(user config.User) (*Socks5ProxyHandler, error) {
	socket, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return nil, err
	}

	port := socket.Addr().(*net.TCPAddr).Port
	fmt.Printf("socks5 proxy server auth [%s] on port [:%d], listening ... \n", aurora.Green(user.User), aurora.Green(port))

	env.SetSocket5Port(user.User, port)

	for {
		client, err := socket.Accept()

		if err != nil {
			log.Errorf("客户端连接失败：%s", err)
			continue
		}

		var handler Handler = &Socks5ProxyHandler{
			Port: port,
			Auth: false,
			User: user.User,
			Pass: user.Pass,
		}

		go handler.Handle(client)

		log.Println(aurora.Blue(client), " request handling...")
	}
}
