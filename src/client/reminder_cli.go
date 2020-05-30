package main

import (
	"MailReminder/src/client/support"
	"flag"
	"log"
	"net"
)

var ipAddr = flag.String("ip", "localhost", "ip address of host")
var port = flag.String("p", "8888", "port number")
var newMsg = flag.Bool("new", false, "upload new reminder to server")
var delMsg = flag.Bool("del", false, "delete an existing reminder")
var lsMsg = flag.Bool("ls", false, "list existing reminders")

func main() {
	flag.Parse()
	checkFlag()

	tcpAddr, err := net.ResolveTCPAddr("tcp4", *ipAddr+":"+*port)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	if *newMsg {
		support.NewMsg(conn)
		return
	}
}

func checkFlag() {
	count := 0
	if *newMsg {
		count++
	}
	if *delMsg {
		count++
	}
	if *lsMsg {
		count++
	}
	if count != 1 {
		log.Fatal("Invalid flag, try -h")
	}
}
