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

var _ Model = &Menu{}

func (m *Menu) Insert(db *sql.DB, params ...interface{}) error {
	query := `INSERT INTO Menu(ID, REST_ID, PRODUCT_ID, PRICE)
			  SELECT Menu.id, Restaurant.id, Product.id, price FROM Restaurant inner join Menu on Restaurant.id = Menu.rest_id inner join Product on Menu.product_id = Product.id
			  WHERE Menu.id = ?`
	_, err := db.Exec(query, m.Id, params[0])
	return err
}
