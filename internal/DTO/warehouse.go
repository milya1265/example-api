package DTO

type ReqGetProducts struct {
	WarehouseID int `json:"warehouse_id"`
}

type ResGetProducts struct {
	ProductsCodes []string `json:"products_codes"`
}
