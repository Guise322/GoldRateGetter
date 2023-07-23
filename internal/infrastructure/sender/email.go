package sender

import (
	"crypto/tls"

	gomail "gopkg.in/mail.v2"
)

const (
	from     = "guise322@ya.ru"
	pass     = "nxwamiqmoqdolhds"
	to       = "dimsonex@ya.ru"
	smtpHost = "smtp.yandex.ru"
	smtpPort = 465
)

type Sender struct{}

func (s Sender) Send(price float32) {
	msg := "Hello there!"
	sub := "testSub"
	m := gomail.NewMessage()
	configure(m, sub, msg)
	d := gomail.NewDialer(smtpHost, smtpPort, from, pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func configure(m *gomail.Message, sub, msg string) {
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", sub)
	m.SetBody("text/plain", msg)
}
