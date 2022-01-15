package models

import "time"

type Restaurant struct {
	Id      int       `json:"id,omitempty"`
	Name    string    `json:"name,omitempty"`
	Image   string    `json:"image,omitempty"`
	OpenAt  time.Time `json:"opening"`
	CloseAt time.Time `json:"closing"`
}
