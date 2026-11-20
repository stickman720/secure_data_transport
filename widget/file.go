package widget

import (
	"io"
	"log"
	"time"

	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)


type MinioHandler struct {
	client *minio.Client
}


type MinioError struct {
	smg string
}

func (err MinioError)Error()string{
	return err.smg
}




func NewMinioHandler(addresss, accesskey, accesssecret string ) (MinioHandler , error ){
	client, err := minio.New(addresss, &minio.Options{
		Creds:  credentials.NewStaticV4(accesskey, accesssecret, ""),
		Secure: false,
	})

	if err != nil{
		log.Println(err.Error())
		return MinioHandler{} , MinioError{smg: "unable to connect to object storage"}
	}

	return MinioHandler{
		client: client,
	} , nil
}


func (mh MinioHandler)NewBucket(ctx context.Context, name , location string) error {
	to , cancel:= context.WithTimeout(ctx , time.Second)
	defer cancel()
	err := mh.client.MakeBucket(to , name , minio.MakeBucketOptions{
		Region: location,
	})

	return err
}



func (mh MinioHandler)UploadFile(ctx context.Context , bucket , oname string , size int64, obj io.Reader ) (minio.UploadInfo , error ){
	to , cancel:= context.WithTimeout(ctx , time.Second )
	defer cancel()
	return mh.client.PutObject(to , bucket , oname , obj , size , minio.PutObjectOptions{})
}



func (mh MinioHandler)DeleteFile(ctx context.Context, bucket , oname string ) error {
	to , cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return mh.client.RemoveObject(to , bucket , oname , minio.RemoveObjectOptions{})
}



func (mh MinioHandler)GetFile(ctx context.Context, bucket, oname string) (io.Reader , error ) {
	return mh.client.GetObject(ctx, bucket, oname, minio.GetObjectOptions{})
}