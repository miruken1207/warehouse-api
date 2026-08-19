package model

type Warehouse struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Location string `json:"location" db:"location"`
}

type CreateWarehouseRequest struct {
	Name     string `json:"name" db:"name"`
	Location string `json:"location" db:"location"`
}
