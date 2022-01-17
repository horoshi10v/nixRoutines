package models

import "database/sql"

type Model interface {
	Insert(db *sql.DB, params ...interface{}) error
}
