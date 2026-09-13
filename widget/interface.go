package widget

import "context"


type Mailer interface {
	SendMail (to []string , subject , body string ) error
}



type Cach interface {
	// exp time is minot
	Set (ctx context.Context , key string , value string , exp_time int) error
	Get (ctx context.Context , key string) (string , error) 
}