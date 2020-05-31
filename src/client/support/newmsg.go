package support

import (
	"MailReminder/src/protocol"
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	New byte = iota
	Del
	Ls
)

const (
	emailPattern = `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	timeFormat   = "2006-01-02 15:04:05"
)

var emailAddrReg = regexp.MustCompile(emailPattern)

func NewMsg(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Send to (e-mail addr separated by spaces): ")
	mailTo, err := readEmailAddr(reader)
	for err != nil {
		fmt.Println(err)
		fmt.Print("Send to (e-mail addr separated by spaces): ")
		mailTo, err = readEmailAddr(reader)
	}

	fmt.Print("Subject: ")
	subject, err := readLine(reader)
	for err != nil {
		fmt.Println(err)
		fmt.Print("Subject: ")
		subject, err = readLine(reader)
	}

	fmt.Print("Body (text/html): ")
	body, err := readLine(reader)
	for err != nil {
		fmt.Println(err)
		fmt.Print("Body (text/html): ")
		body, err = readLine(reader)
	}

	fmt.Print("Time (yyyy-MM-dd HH:mm:ss): ")
	sendTime, err := readTime(reader)
	for err != nil {
		fmt.Println(err)
		fmt.Print("Time (yyyy-MM-dd HH:mm:ss): ")
		sendTime, err = readTime(reader)
	}

	datagram := &protocol.Datagram{
		Op:       New,
		Subject:  subject,
		Body:     body,
		SendTime: sendTime,
		SendTo:   mailTo,
	}
	jsonBytes, err := datagram.GetJsonBytes()
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Write(jsonBytes)
	if err != nil {
		log.Fatal(err)
	}

	resp := make([]byte, 1)
	_, err = conn.Read(resp)
	if err != nil {
		log.Fatal(err)
	}
	if resp[0] == 0 {
		fmt.Println("success")
		return
	}
	fmt.Printf("subject duplicate: '%s', try another subject\n", subject)
}

func readEmailAddr(reader *bufio.Reader) ([]string, error) {
	input, err := readLine(reader)
	if err != nil {
		return nil, err
	}

	mailTo := strings.Split(input, " ")
	for _, addr := range mailTo {
		if !emailAddrReg.MatchString(addr) {
			return nil, errors.New("illegal email address: " + addr)
		}
	}
	return mailTo, nil
}

func readLine(reader *bufio.Reader) (string, error) {
	str, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	str = strings.ReplaceAll(str, "\r\n", "")
	str = strings.ReplaceAll(str, "\n", "")
	return str, nil
}

func readTime(reader *bufio.Reader) (string, error) {
	input, err := readLine(reader)
	if err != nil {
		return "", err
	}
	_, err = time.Parse(timeFormat, input)
	if err != nil {
		return "", err
	}
	return input, nil
}
