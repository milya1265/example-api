package model

type Product struct {
	ID         int    `json:"id"`
	UniqueCode string `json:"unique_code"`
	Name       string `json:"name"`
	Size       int    `json:"size"`
	Count      int    `json:"count"`
	Left       int    `json:"left"`
}

type Reservation struct {
	ID          int    `json:"id"`
	WarehouseID int    `json:"warehouse_id"`
	ProductCode string `json:"product_id"`
	Count       int    `json:"count"`
}
