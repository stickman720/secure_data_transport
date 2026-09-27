package main 

import (
	"context"
	_"embed"
	"os"
	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5"
)


//go:embed sqlc/schema.sql
var ddl string

func migrate(){

	ctx := context.Background()


	err := godotenv.Load()
	if err != nil {
		panic(err)
	}



	db, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close(ctx)

	

}