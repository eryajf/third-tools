package main

import (
	"fmt"

	"strconv"

	"gopkg.in/gomail.v2"
)

var (
	EmailUser string = "Linuxlql@163.com"
	EmailPass string = "xxxx"
	EmailHost string = "smtp.163.com"
	EmailPort string = "465"
	EmailForm string = "go-ldap-admin后台"
)

func main() {
	// 接收两个参数，第一个是要接收邮件的邮箱，第二个是新的密码
	SendMail([]string{"eryajf@163.com"}, "testPass")
}

func email(mailTo []string, subject string, body string) error {
	mailConn := map[string]string{
		"user": EmailUser,
		"pass": EmailPass,
		"host": EmailHost,
		"port": EmailPort,
	}
	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	newmail := gomail.NewMessage()

	newmail.SetHeader("From", newmail.FormatAddress(mailConn["user"], EmailForm))
	newmail.SetHeader("To", mailTo...)    //发送给多个用户
	newmail.SetHeader("Subject", subject) //设置邮件主题
	newmail.SetBody("text/html", body)    //设置邮件正文

	do := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])
	return do.DialAndSend(newmail)
}

func SendMail(sendto []string, pass string) error {
	subject := "重置LDAP密码成功"
	// 邮件正文
	body := fmt.Sprintf("<li><a>更改之后的密码为:%s</a></li>", pass)
	return email(sendto, subject, body)
}
