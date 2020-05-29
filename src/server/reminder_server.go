package main

import (
	"MailReminder/src/conf"
	"MailReminder/src/reminder"
	"log"
	"net"
	"os"
)

var separator = string(os.PathSeparator)
var confPath = ".." + separator + "resource" + separator + "config.xml"

func main() {
	config, err := conf.GetConfig(confPath)
	if err != nil {
		log.Fatal(err)
	}

	mailReminder, err := reminder.NewReminder(config)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", "0.0.0.0:"+config.Listen)
	if err != nil {
		log.Fatal(err)
	}

	go mailReminder.Serve()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn, mailReminder)
	}
}
