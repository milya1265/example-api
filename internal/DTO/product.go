package DTO

type ReqReserveProduct struct {
	WarehouseID int      `json:"warehouse_id"`
	UniqueCodes []string `json:"unique_codes"`
	Counts      []int    `json:"counts"`
}

type ResReserveProduct struct {
	Successful   []Successful `json:"successful"`
	Unsuccessful []string     `json:"unsuccessful"`
	Errors       []string     `json:"errors"`
}

type Successful struct {
	ID         int    `json:"id"`
	UniqueCode string `json:"unique_codes"`
}

type ReqFreeReservation struct {
	ID []int `json:"id"`
}

type ResFreeReservation struct {
	Successful   []int    `json:"successful"`
	Unsuccessful []int    `json:"unsuccessful"`
	Errors       []string `json:"errors"`
}
