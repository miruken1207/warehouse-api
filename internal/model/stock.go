package model

type Stock struct {
	ID          int `json:"id" db:"id"`
	WarehouseID int `json:"warehouse_id" db:"warehouse_id"`
	ItemID      int `json:"item_id" db:"item_id"`
	Quantity    int `json:"quantity" db:"quantity"`
}

type CreateStockRequest struct {
	WarehouseID int `json:"warehouse_id"`
	ItemID      int `json:"item_id"`
	Quantity    int `json:"quantity"`
}

type UpdateStockRequest struct {
	WarehouseID int `json:"warehouse_id"`
	ItemID      int `json:"item_id"`
	Delta       int `json:"delta"`
}

type TransferStockRequest struct {
	FromWarehouseID int `json:"from_warehouse_id"`
	ToWarehouseID   int `json:"to_warehouse_id"`
	ItemID          int `json:"item_id"`
	Quantity        int `json:"quantity"`
}
