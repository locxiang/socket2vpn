package main

import (
	"github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gopkg.in/ffmt.v1"
	"net"
	"os"
	"path"
	"runtime"
	"socket2vpn/config"
	"socket2vpn/proxy"
)

var (
	cfg = pflag.StringP("config", "c", "", "config file path.")
)

func main() {

	//使用全部cpu
	runtime.GOMAXPROCS(runtime.NumCPU())
	pflag.Parse()

	// 初始化配置文件
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	ffmt.Print(config.Values)
	initLog()

	socket, err := net.Listen("tcp", ":8081")
	if err != nil {
		return
	}
	log.Infof("socks5 proxy server running on port [:%d], listening ...", aurora.Green(8081))

	for {
		client, err := socket.Accept()

		if err != nil {
			return
		}

		var handler proxy.Handler = &proxy.Socks5ProxyHandler{
			Auth: false,
		}

		go handler.Handle(client)

		log.Println(aurora.Blue(client), " request handling...")
	}

}

func initLog() {
	c := config.Values.Log

	if c.LoggerFile != "" && c.Writers == "file" {
		_ = os.MkdirAll(path.Dir(c.LoggerFile), os.ModePerm)
		file, err := os.OpenFile(c.LoggerFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
		if err != nil {
			log.Panicf("log  init failed:%s", err)
		}

		log.SetFormatter(&log.TextFormatter{
			DisableColors: true,
		})

		log.SetOutput(file)
	} else {
		log.SetOutput(os.Stdout)
	}

	if c.LoggerLevel == "ERROR" {
		log.SetLevel(log.ErrorLevel)
	}
}
