package proxy

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestSocks5ProxyHandler_getPPPIp(t *testing.T) {

	s := Socks5ProxyHandler{Auth: false}
	ip, _ := s.getPPPIp("ppp0")
	fmt.Printf("%s", ip)
}

func BenchmarkSocks5ProxyHandler_Handle(b *testing.B) {
	s := Socks5ProxyHandler{Auth: false}

	for i := 0; i < b.N; i++ {
		s.getPPPIp("ppp0")
	}
}

func TestXX(t *testing.T) {
	s := "ppp0"
	p := unsafe.Pointer(&[]byte(s)[0])
	fmt.Printf("%+v", p)
}
