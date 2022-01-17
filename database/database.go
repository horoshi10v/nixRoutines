package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	dbConn, err := sql.Open(
		"mysql",
		"Khoroshylov:Valentyn!1@/routines_db?parseTime=true",
	)
	if err != nil {
		return nil, err
	}

	err = dbConn.Ping()
	if err != nil {
		return nil, err
	}
	return dbConn, err
}
func DropTable1(db *sql.DB) {
	query := "DELETE FROM restaurant"
	db.Exec(query)
}
func DropTable2(db *sql.DB) {
	query := "DELETE FROM menu"
	db.Exec(query)
}
