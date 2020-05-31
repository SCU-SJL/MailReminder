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
		newMail(datagram, rd, conn)
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

func newMail(datagram *protocol.Datagram, rd *reminder.Reminder, conn net.Conn) {
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
	ok := rd.NewMsg(mail)
	if ok {
		_, _ = conn.Write([]byte{0})
		return
	}
	_, _ = conn.Write([]byte{1})
}

func listMails(rd *reminder.Reminder, conn net.Conn) {
	mails := struct {
		Subjects  []string   `json:"subjects"`
		Receivers [][]string `json:"receivers"`
		SendTime  []string   `json:"send_time"`
	}{}
	mails.Subjects, mails.Receivers, mails.SendTime = rd.GetMailList()

	jsonBytes, err := json.Marshal(mails)
	if err != nil {
		return
	}
	_, _ = conn.Write(jsonBytes)
}

func delMail(datagram *protocol.Datagram, rd *reminder.Reminder, conn net.Conn) {
	status := rd.DelMsg(datagram.Subject)
	_, _ = conn.Write([]byte{status})
}
