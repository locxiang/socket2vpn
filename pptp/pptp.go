package pptp

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os/exec"
)

// 简历pptp通道
func Conn() {

	fmt.Printf("start ")
	cmd := exec.Command("ifconfig")

	stdout, err := cmd.StdoutPipe()
	if err != nil { //获取输出对象，可以从该对象中读取输出结果
		log.Fatal(err)
	}
	defer stdout.Close() // 保证关闭输出流

	cmd.Start()

	if opBytes, err := ioutil.ReadAll(stdout); err != nil { // 读取输出结果
		log.Fatal(err)
	} else {
		fmt.Printf("文件： %s", opBytes)
	}

	cmd.Wait()

}
