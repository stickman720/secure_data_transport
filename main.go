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
	mailer := widget.Mailer{
		From:     os.Getenv("EMAIL"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("PORT"),
	}
	_ = mailer

	db, err := sql.Open("sqlite","./secure_data_transport.db")
	if err != nil {
		panic(err)
	}



	queries := sqlcpkg.New(db)


	ag := r.Group("/api/auth")
	api.InitAuthRoutes(ag , queries)

	


	r.Run(":8000")

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