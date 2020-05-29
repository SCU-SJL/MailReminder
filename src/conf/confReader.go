package conf

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"strconv"
)

type ReminderConfig struct {
	XmlName xml.Name `xml:"config"`
	Addr    string   `xml:"addr"`
	Auth    string   `xml:"auth"`
	Host    string   `xml:"host"`
	Port    string   `xml:"port"`
	Retry   string   `xml:"retry"`
	Max     string   `xml:"max"`
	Listen  string   `xml:"listen"`
}

func (conf *ReminderConfig) GetMax() (int, error) {
	return strconv.Atoi(conf.Max)
}

func (conf *ReminderConfig) GetHost() string {
	return conf.Host
}

func (conf *ReminderConfig) GetRetry() (int, error) {
	return strconv.Atoi(conf.Retry)
}

func (conf *ReminderConfig) GetPort() (int, error) {
	return strconv.Atoi(conf.Port)
}

func (conf *ReminderConfig) GetAuth() string {
	return conf.Auth
}

func (conf *ReminderConfig) GetAddr() string {
	return conf.Addr
}

func (conf *ReminderConfig) GetListen() (int, error) {
	return strconv.Atoi(conf.Listen)
}

func GetConfig(path string) (*ReminderConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	conf := ReminderConfig{}
	err = xml.Unmarshal(content, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
