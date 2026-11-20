package services

import (
	//"context"
	//"io"
	"secure_data_transport/widget"
	//"github.com/minio/minio-go/v7"
	sqlcpkg "secure_data_transport/sqlcpkg"
	"github.com/jackc/pgx/v5"


)


type FileService struct {
	filehandler *widget.FileHandler
	queries *sqlcpkg.Queries
	conn *pgx.Conn
}



func NewFileHandler(filehandler *widget.FileHandler , queries *sqlcpkg.Queries , conn *pgx.Conn) FileService{
	return FileService{
		filehandler: filehandler,
		queries: queries,
		conn: conn,
	}
}

