package main

import (
	"database/sql"
	"os"
	"secure_data_transport/api"
	sqlcpkg "secure_data_transport/sqlcpkg"
	"secure_data_transport/widget"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
	"github.com/joho/godotenv"
)





func startapp(){

		r := gin.Default()

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

	
	

	var cach widget.Cach = widget.NewRedisCach(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"))
	



	db, err := sql.Open("sqlite","./secure_data_transport.db")
	if err != nil {
		panic(err)
	}



	queries := sqlcpkg.New(db)


	ag := r.Group("/api/auth")
	api.InitAuthRoutes(ag , queries , &cach , mailer)

	


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