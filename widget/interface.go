package widget

import (
	"context"
	"time"
)


type Mailer interface {
	SendMail (to []string , subject , body string ) error
}



type Cach interface {
	// exp time is minot
	Set (ctx context.Context , key string , value interface{} , exp_time time.Duration) error
	Get (ctx context.Context , key string , out any) (error) 
}