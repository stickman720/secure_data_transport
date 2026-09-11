package main 

import (
	"context"
	_"embed"
	"database/sql"
)


//go:embed sqlc/schema.sql
var ddl string

func migrate(){

	ctx := context.Background()

	db, err := sql.Open("sqlite","./secure_data_transport.db")
	if err != nil {
		panic(err)
	}

	if _ , err := db.ExecContext(ctx, ddl); err != nil {
		panic(err)
	}

}