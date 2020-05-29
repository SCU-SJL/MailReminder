package main

import (
	"MailReminder/src/reminder"
	"gopkg.in/gomail.v2"
	"net"
	"strings"
	"time"
)

const (
	New byte = iota
	Del
	Ls
	Exit
)

const format = "2006-01-02 15:04:05"

func handleConn(conn net.Conn, rd *reminder.Reminder) {
	op, err := readOp(conn)
	if err != nil {
		return
	}
	switch op {
	case New:
		newMail(conn, rd)
	}
}

func readOp(conn net.Conn) (byte, error) {
	buf := make([]byte, 1)
	_, err := conn.Read(buf)
	return buf[0], err
}

func newMail(conn net.Conn, rd *reminder.Reminder) {
	mailTo, err := readMailTo(conn)
	if err != nil {
		return
	}
	subject := "Mail-Reminder rings!"
	body, err := readBody(conn)
	if err != nil {
		return
	}
	sendTime, err := readBody(conn)
	if err != nil {
		return
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", msg.FormatAddress(rd.User, "Mail-Reminder"))
	msg.SetHeader("To", mailTo...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	t, _ := time.Parse(format, sendTime)
	mail := &reminder.Mail{
		Msg:  msg,
		Time: t,
	}
	rd.NewMsg(mail)
}

func readMailTo(conn net.Conn) ([]string, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	bufStr := string(buf[:n])
	return strings.Split(bufStr, " "), nil
}

func readBody(conn net.Conn) (string, error) {
	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}
