package main

import (
	"MailReminder/src/conf"
	"MailReminder/src/reminder"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"time"
)

var separator = string(os.PathSeparator)

func main() {
	config, err := conf.GetConfig(".." + separator + "resource" + separator + "config.xml")
	if err != nil {
		log.Fatal(err)
	}

	mailReminder, err := reminder.NewReminder(config)
	if err != nil {
		log.Fatal(err)
	}

	mailTo := []string{
		"scu_sjl@outlook.com",
		"sjl66666666@gmail.com",
	}
	subject := "MailReminder测试"
	body := "Hello, MailReminder向您问好"

	msg := gomail.NewMessage()
	msg.SetHeader("From", msg.FormatAddress(config.GetAddr(), "MailReminder"))
	msg.SetHeader("To", mailTo...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)
	format := "2006-01-02 15:04:05"
	t, _ := time.Parse(format, "2020-05-28 16:05:20")
	mail := &reminder.Mail{
		Msg:  msg,
		Time: t,
	}
	go mailReminder.Serve()
	mailReminder.NewMsg(mail)
	<-time.After(30 * time.Second)
}
