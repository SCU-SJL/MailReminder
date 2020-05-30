package protocol

import "encoding/json"

type Datagram struct {
	Op       byte     `json:"op"`
	Id       int      `json:"id"`
	Subject  string   `json:"subject"`
	Body     string   `json:"body"`
	SendTime string   `json:"send_time"`
	SendTo   []string `json:"send_to"`
}

func (data *Datagram) GetJsonBytes() ([]byte, error) {
	return json.Marshal(*data)
}

func ConvertToDatagram(byteArr []byte) (*Datagram, error) {
	var datagram Datagram
	err := json.Unmarshal(byteArr, &datagram)
	if err != nil {
		return nil, err
	}
	return &datagram, nil
}
