package widget

import (
	"net/smtp"
	"log"
)



type SMTPMailer struct {
	From     string
	Password string
	Host     string
	Port     string
}



func (m SMTPMailer) SendMail(to []string , subject , body string) error {
	from := m.From
	password := m.Password



	smtpHost := m.Host
	smtpPort := m.Port

	message := []byte("subject : " + subject + "\r\n\r\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	_ = auth
	err := smtp.SendMail(smtpHost+":"+smtpPort, nil, from, to, message)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}