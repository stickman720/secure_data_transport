package widget

import (
	"fmt"
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
	smtpPort := "587"

	message := []byte("subject : " + subject + "\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Email sent successfully!")

	return nil
}