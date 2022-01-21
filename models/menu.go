package models

import "database/sql"

type Menu struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	Image       string   `json:"image"`
	Type        string   `json:"type"`
	Ingredients []string `json:"ingredients"`
}

func (prod *Menu) Insert(db *sql.DB) (int64, error) {
	res, err := db.Exec("INSERT INTO product VALUE (?, ?, ?, ?, ?)",
		prod.Id, prod.Name, prod.Price, prod.Image, prod.Type)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
