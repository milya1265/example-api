package model

type Warehouse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Availability bool   `json:"availability"`
}
