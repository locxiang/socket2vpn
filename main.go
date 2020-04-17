package main

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gopkg.in/ffmt.v1"
	"os"
	"path"
	"runtime"
	"socket2vpn/config"
	"socket2vpn/pptp"
	"socket2vpn/proxy"
	"time"
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

	users := config.Values.Users

	pptp.CloseAll()

	for _, u := range users {
		if err := pptp.NewPPTP(u); err != nil {
			log.Fatalf("建立[%s]pptp通道出错: %s", aurora.BgRed(u.User), err)
		}
		go proxy.NewSocket5(u)
	}

	select {}
}

func initLog() {

	time.Sleep(1 * time.Second)
	c := config.Values.Log

	log.SetFormatter(&nested.Formatter{
		TimestampFormat: time.RFC3339,
	})

	if c.LoggerFile != "" && c.Writers == "file" {
		_ = os.MkdirAll(path.Dir(c.LoggerFile), os.ModePerm)
		file, err := os.OpenFile(c.LoggerFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
		if err != nil {
			log.Panicf("log  init failed:%s", err)
		}

		log.SetOutput(file)
	} else {
		log.SetOutput(os.Stdout)
	}

	if c.LoggerLevel == "ERROR" {
		log.SetLevel(log.ErrorLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}
}
