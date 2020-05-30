package support

import (
	"MailReminder/src/protocol"
	"MailReminder/src/reminder"
	"encoding/json"
	"gopkg.in/gomail.v2"
	"net"
	"time"
)

const (
	New byte = iota
	Del
	Ls
)

const format = "2006-01-02 15:04:05"

func HandleConn(conn net.Conn, rd *reminder.Reminder) {
	datagram, err := readDatagram(conn)
	if err != nil {
		return
	}
	switch datagram.Op {
	case New:
		newMail(datagram, rd)
	case Del:
		delMail(datagram, rd, conn)
	case Ls:
		listMails(rd, conn)
	}
}

func readDatagram(conn net.Conn) (*protocol.Datagram, error) {
	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return protocol.ConvertToDatagram(buf[:n])
}

func newMail(datagram *protocol.Datagram, rd *reminder.Reminder) {
	msg := gomail.NewMessage()
	msg.SetHeader("From", msg.FormatAddress(rd.User, "Mail-Reminder"))
	msg.SetHeader("To", datagram.SendTo...)
	msg.SetHeader("Subject", datagram.Subject)
	msg.SetBody("text/html", datagram.Body)

	t, _ := time.Parse(format, datagram.SendTime)
	mail := &reminder.Mail{
		Msg:  msg,
		Time: t,
	}
	rd.NewMsg(mail)
}

func listMails(rd *reminder.Reminder, conn net.Conn) {
	subjects := struct {
		Data []string `json:"data"`
	}{Data: rd.GetAllSubjects()}

	jsonBytes, err := json.Marshal(subjects)
	if err != nil {
		return
	}
	_, _ = conn.Write(jsonBytes)
}

func delMail(datagram *protocol.Datagram, rd *reminder.Reminder, conn net.Conn) {
	ok := rd.DelMsg(datagram.Id)
	if ok {
		_, _ = conn.Write([]byte{0})
		return
	}
	_, _ = conn.Write([]byte{1})
}
