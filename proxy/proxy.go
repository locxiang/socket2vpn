package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	//"syscall"
)

var (
	noAuth   = []byte{0x05, 0x00}
	withAuth = []byte{0x05, 0x02}

	authSuccess = []byte{0x05, 0x00}
	authFailed  = []byte{0x05, 0x01}

	connectSuccess = []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

type Socks5ProxyHandler struct {
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

func (Socks5ProxyHandler) getPPPIp(netInterface string) (*net.TCPAddr, error) {
	ief, err := net.InterfaceByName(netInterface)
	if err != nil {
		fmt.Println(err)
	}
	addrs, err := ief.Addrs()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	tcpAddr := &net.TCPAddr{
		IP: addrs[0].(*net.IPNet).IP,
	}

	return tcpAddr, nil
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

	if b[0] == 0x05 {

		if socks5.Auth == false {
			connect.Write(noAuth)
		} else {
			connect.Write(withAuth)

			_, err = connect.Read(b)
			if err != nil {
				return
			}

			user_length := int(b[1])
			user := string(b[2:(2 + user_length)])
			pass := string(b[(2 + user_length):])

			if socks5.User == user && socks5.Pass == pass {
				connect.Write(authSuccess)
			} else {
				connect.Write(authFailed)
				return
			}
		}

		addr, err := socks5.getPPPIp("ppp0")
		if err != nil {
			log.Printf("获取ip失败")
			return
		}

		dialer := net.Dialer{
			LocalAddr: addr,
			//Control: func(network, address string, c syscall.RawConn) error {
			//
			//	log.Printf("net:%s add:%s",network,address)
			//	return c.Control(func(fd uintptr) {
			//		err := syscall.SetsockoptString(int(fd), syscall.SOL_SOCKET, 25, "ppp0")
			//		if err != nil {
			//			log.Printf("control: %s", err)
			//			return
			//		}
			//	})
			//},
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
}
