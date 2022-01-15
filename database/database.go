package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func Connect() (*sql.DB, error) {
	dbConn, err := sql.Open(
		"mysql",
		"Khoroshylov:Valentyn!1@/routines_db?parseTime=true",
	)
	if err != nil {
		log.Fatal(err)
	}
	return dbConn, err
}
