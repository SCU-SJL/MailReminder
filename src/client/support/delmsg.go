package support

import (
	"MailReminder/src/protocol"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func DelMsg(conn net.Conn) {
	r := bufio.NewReader(os.Stdin)
	prefix := readPrefix(r)
	datagram := &protocol.Datagram{
		Op:      Del,
		Subject: prefix,
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
	if resp[0] == 1 {
		fmt.Printf("There is no mail has prefix '%s', try -ls to check mail list", prefix)
		return
	}
	fmt.Printf("More than one mail has prefix '%s', try to enter a longer one\n", prefix)
}

func readPrefix(r *bufio.Reader) string {
	fmt.Print("Enter the prefix of mail subject: ")
	prefix, err := r.ReadString('\n')
	for err != nil {
		fmt.Print("Enter the prefix of mail subject: ")
		prefix, err = r.ReadString('\n')
	}
	prefix = strings.ReplaceAll(prefix, "\r\n", "")
	prefix = strings.ReplaceAll(prefix, "\n", "")
	return prefix
}
