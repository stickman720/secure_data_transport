package main

import (
	"context"
	"os"
	"secure_data_transport/api"
	"secure_data_transport/core"
	sqlcpkg "secure_data_transport/sqlcpkg"
	"secure_data_transport/widget"
    "github.com/jackc/pgx/v5"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)





func startapp(){

	r := gin.Default()
	ctx := context.Background()


	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	
	var mailer widget.Mailer = widget.SMTPMailer{
		From:     os.Getenv("EMAIL"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("PORT"),
	}
	//==============================================================
	var cach widget.Cach = widget.NewRedisCach(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"))
	//==============================================================
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	queries := sqlcpkg.New(conn)
	//==============================================================
	security := core.Security{
		JWTSecret: os.Getenv("JWT_SECRET"),
		HashSecret: os.Getenv("HASH_SECRET"),
	}



	ag := r.Group("/api/auth")
	api.InitAuthRoutes(ag , queries , &cach , mailer , security)

	


	r.Run(":8080")

}







func main() {

	args := os.Args

	switch args[1]{
	case "startapp":
		startapp()
	case "migrate":
		migrate()
	}
	

}