package util

import (
	"fmt"
	"net"
)

func GetPPPIp(connectName string) (*net.TCPAddr, error) {
	ief, err := net.InterfaceByName(connectName)
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
