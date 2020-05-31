package reminder

import (
	"MailReminder/src/conf"
	"gopkg.in/gomail.v2"
	"log"
	"mime"
	"strings"
	"sync"
	"time"
)

var decoder mime.WordDecoder

type Reminder struct {
	User  string
	auth  string
	host  string
	port  int
	retry int
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

func (reminder *Reminder) Serve() {
	log.Printf("[start] mail reminder\n")
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

func (reminder *Reminder) NewMsg(msg *Mail) bool {
	reminder.mu.Lock()
	defer reminder.mu.Unlock()
	subject := msg.Msg.GetHeader("Subject")[0]
	for _, m := range reminder.mails {
		if m.Msg.GetHeader("Subject")[0] == subject {
			return false
		}
	}
	reminder.mails = append(reminder.mails, msg)
	title, err := decoder.Decode(subject)
	if err != nil {
		title = subject
	}
	log.Printf("[new] subject: %s, sendTime: %s\n", title, msg.Time.Format("2006-01-02 15:04:05"))
	return true
}

func (reminder *Reminder) GetMailList() (subjects []string, receivers [][]string, sendTime []string) {
	reminder.mu.Lock()
	defer reminder.mu.Unlock()
	for _, m := range reminder.mails {
		subject, err := decoder.Decode(m.Msg.GetHeader("Subject")[0])
		if err != nil {
			subject = m.Msg.GetHeader("Subject")[0]
		}
		subjects = append(subjects, subject)
		receivers = append(receivers, m.Msg.GetHeader("To"))
		sendTime = append(sendTime, m.Time.Format("2006-01-02 15:04:05"))
	}
	return
}

func (reminder *Reminder) DelMsg(prefix string) byte {
	reminder.mu.Lock()
	defer reminder.mu.Unlock()
	var target = -1
	var targetSub string
	for i, m := range reminder.mails {
		subject, _ := decoder.Decode(m.Msg.GetHeader("Subject")[0])
		if strings.HasPrefix(subject, prefix) {
			if target >= 0 {
				return 2
			}
			target = i
			targetSub = subject
		}
	}
	if target < 0 {
		return 1
	}
	reminder.mails = append(reminder.mails[:target], reminder.mails[target+1:]...)
	log.Printf("[del] subject: %s\n", targetSub)
	return 0
}

func (reminder *Reminder) sendMail(msg *gomail.Message) error {
	sub, _ := decoder.Decode(msg.GetHeader("Subject")[0])
	log.Printf("[sending] subject: %s\n", sub)

	d := gomail.NewDialer(reminder.host, reminder.port, reminder.User, reminder.auth)
	sendErr := d.DialAndSend(msg)
	if sendErr != nil {
		for i := 0; i < reminder.retry; i++ {
			if sendErr = d.DialAndSend(msg); sendErr == nil {
				return nil
			}
		}
	}
	if sendErr != nil {
		log.Printf("[failed] subject: %s\n", sub)
	} else {
		log.Printf("[success] subject: %s\n", sub)
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
		User:  conf.GetAddr(),
		auth:  conf.GetAuth(),
		host:  conf.GetHost(),
		port:  p,
		retry: r,
		mails: make([]*Mail, 0, m),
	}, nil
}
