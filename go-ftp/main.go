package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

func main() {
	// 通过ftp内置的Dial连接远程ftp,获得一个连接对象c
	c, err := ftp.Dial("192.168.0.22:21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		fmt.Printf("conn ftp failed, err:%v\n", err)
		return
	}
	// Login 使用用户名密码进行认证
	err = c.Login("ftp_usera", "123456")
	if err != nil {
		fmt.Printf("login ftp failed, err:%v\n", err)
		return
	}
	// 此处定义一个命令行参数，以定义将要下载的文件，如果文件不在根目录，可以使用全路径
	name := flag.String("file", "test-file.txt", "请输入将要下载的文件路径")
	flag.Parse() // 解析命令行参数，千万不要忘了这个参数
	// 创建一个读取文件内容的对象
	r, err := c.Retr(*name)
	if err != nil {
		fmt.Println("retr file failed, err", err)
		return
	}
	// 使用ioutil读取刚刚对象的内容
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Printf("read file failed, err:%v\n", err)
		return
	}
	// 定义文件下载之后保存在本地的路径，因为这里是一个固定的位置，因此写死了，在固定路径下，按天进行分割保存
	path := "/data/www/storage/jzbbankcode" + "/" + time.Now().Format("20060102")
	err = os.MkdirAll(path, 0755) //先创建如上定义的路径
	if err != nil {
		fmt.Printf("mkdir directory failed, err:%v\n", err)
		return
	}
	// 通过截取定义用户输入的文件路径最后一段，从而获取到文件名，以放入本地
	s1 := strings.Split(*name, "/")
	_file := path + "/" + s1[len(s1)-1]
	err = ioutil.WriteFile(_file, []byte(buf), 0644)
	if err != nil {
		fmt.Printf("write file failed, err:%v\n", err)
		return
	}
}
