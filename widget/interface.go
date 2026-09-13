package widget


type Mailer interface {
	SendMail (to []string , subject , body string ) error
}



