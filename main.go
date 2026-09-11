package main

import (
	"database/sql"
	"os"
	"secure_data_transport/api"
	sqlcpkg "secure_data_transport/sqlcpkg"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)





func startapp(){

		r := gin.Default()

	

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