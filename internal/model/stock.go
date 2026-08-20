package model

type Stock struct {
	ID          int `json:"id" db:"id"`
	WarehouseID int `json:"warehouse_id" db:"warehouse_id"`
	ItemID      int `json:"item_id" db:"item_id"`
	Quantity    int `json:"quantity" db:"quantity"`
}

type UpdateStockRequest struct {
	WarehouseID int `json:"warehouse_id"`
	ItemID      int `json:"item_id"`
	Delta       int `json:"delta"`
}
