package models

import "database/sql"

type Restaurant struct {
	Id           int          `json:"id"`
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	Image        string       `json:"image"`
	WorkingHours WorkingHours `json:"workingHours"`
	Menu         []Menu       `json:"menu"`
}

var _ = &Restaurant{}

func (rest *Restaurant) Insert(db *sql.DB) (int64, error) {
	query := `INSERT INTO restaurant(id, name, type, image, open_at, close_at) 
			  VALUES (?,?,?,?,?,?)`
	res, err := db.Exec(query, rest.Id, rest.Name, rest.Type, rest.Image,
		rest.WorkingHours.Opening, rest.WorkingHours.Closing)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
