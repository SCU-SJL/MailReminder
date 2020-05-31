package support

import (
	"MailReminder/src/protocol"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func ListMails(conn net.Conn) {
	datagram := &protocol.Datagram{
		Op: Ls,
	}
	jsonBytes, err := datagram.GetJsonBytes()
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Write(jsonBytes)
	if err != nil {
		log.Fatal(err)
	}

	mails := struct {
		Subjects  []string   `json:"subjects"`
		Receivers [][]string `json:"receivers"`
		SendTime  []string   `json:"send_time"`
	}{}

	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(buf[:n], &mails)
	if err != nil {
		fmt.Println("unmarshal data pkg from server failed")
		log.Fatal(err)
	}

	fmt.Println("[Mail List]:")
	for i := 0; i < len(mails.Subjects); i++ {
		recvStr := fmtReceivers(mails.Receivers[i])

		fmt.Printf("#%2.2d - Subject:\"%s\"\n    - Receivers:\"%s\"\n    - SendTime:\"%s\"\n\n",
			i+1, mails.Subjects[i], recvStr, mails.SendTime[i])
	}
}

func fmtReceivers(receivers []string) string {
	str := ""
	for _, r := range receivers {
		str = str + r + ", "
	}
	if str != "" {
		str = str[:len(str)-2]
	}
	return str
}
