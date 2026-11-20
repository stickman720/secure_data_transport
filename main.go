package main

import (
	"context"
	"secure_data_transport/api"
	"secure_data_transport/config"
	"secure_data_transport/core"
	sqlcpkg "secure_data_transport/sqlcpkg"
	"secure_data_transport/widget"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	_ "modernc.org/sqlite"
)





func startapp(){

	r := gin.Default()
	ctx := context.Background()


	conf , err := config.Init()
	if err != nil{
		panic(err)
	}
	
	var mailer widget.Mailer = widget.SMTPMailer{
		From:     conf.Email,
		Password: conf.EmailPassword,
		Host:     conf.SMTPHost,
		Port:     conf.SMTPPort,
	}
	//==============================================================
	var cach widget.Cach = widget.NewRedisCach(
		conf.RedisAddr,
		conf.RedisPassword,
	)
	//==============================================================
	conn, err := pgx.Connect(ctx, conf.DataBaseURL)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	queries := sqlcpkg.New(conn)
	//==============================================================
	security := core.Security{
		JWTSecret: conf.JWTSecret,
		HashSecret: conf.HashSecret,
	}



	ag := r.Group("/api/auth")
	api.InitAuthRoutes(ag , queries , &cach , mailer , security , conf)
	tg := r.Group("/api/teams")
	api.InitTeamRoutes(tg , queries , &cach , mailer , security , conf)
	

	var filehandler widget.FileHandler 
	filehandler , err = widget.NewMinioHandler("http://minio:9000" , "sdt_user" , "sdt_dev_password")

	
	_ = filehandler


	r.Run(":8080")

}







func main() {


	startapp()
	

}