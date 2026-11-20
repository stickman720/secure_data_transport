package widget

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
)


type Mailer interface {
	SendMail (to []string , subject , body string ) error
}



type Cach interface {
	// exp time is minot
	Set (ctx context.Context , key string , value interface{} , exp_time time.Duration) error
	Get (ctx context.Context , key string , out any) (error) 
	DropKey(ctx context.Context, key string) error
}



type FileHandler interface {
	NewBucket (ctx context.Context, name , location string) error
	UploadFile (ctx context.Context, bucket, name string, size int64, obj io.Reader) (minio.UploadInfo, error)
	GetFile (ctx context.Context, bucket, oname string) (io.Reader , error )
	
}