package reminder

import (
	"MailReminder/src/conf"
	"fmt"
	"gopkg.in/gomail.v2"
	"net"
	"sync"
	"time"
)

type Reminder struct {
	user  string
	auth  string
	host  string
	port  int
	retry int
	conn  net.Conn

	mails []*Mail
	mu    sync.Mutex
}

type Mail struct {
	Msg  *gomail.Message
	Time time.Time
}

func (m *Mail) isReady() bool {
	now, _ := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	return now.After(m.Time)
}

func (reminder *Reminder) NewMsg(msg *Mail) {
	reminder.mu.Lock()
	defer reminder.mu.Unlock()
	reminder.mails = append(reminder.mails, msg)
}

func (reminder *Reminder) Serve() {
	for {
		reminder.mu.Lock()
		for i, mail := range reminder.mails {
			if mail.isReady() {
				sendErr := reminder.sendMail(mail.Msg)
				if sendErr == nil {
					reminder.mails = append(reminder.mails[:i], reminder.mails[i+1:]...)
					break
				}
			}
		}
		reminder.mu.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func (reminder *Reminder) sendMail(msg *gomail.Message) error {
	fmt.Println("sending")
	d := gomail.NewDialer(reminder.host, reminder.port, reminder.user, reminder.auth)
	sendErr := d.DialAndSend(msg)
	if sendErr != nil {
		for i := 0; i < reminder.retry; i++ {
			if sendErr = d.DialAndSend(msg); sendErr == nil {
				return nil
			}
		}
	}
	return sendErr
}

func NewReminder(conf *conf.ReminderConfig) (*Reminder, error) {
	p, err := conf.GetPort()
	if err != nil {
		return nil, err
	}

	r, err := conf.GetRetry()
	if err != nil {
		return nil, err
	}

	m, err := conf.GetMax()
	if err != nil {
		return nil, err
	}

	return &Reminder{
		user:  conf.GetAddr(),
		auth:  conf.GetAuth(),
		host:  conf.GetHost(),
		port:  p,
		retry: r,
		conn:  nil,
		mails: make([]*Mail, 0, m),
	}, nil
}
