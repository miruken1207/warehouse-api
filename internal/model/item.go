package model

type Item struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Category string `json:"category" db:"category"`
}

type CreateItemRequest struct {
	Name     string `json:"name" db:"name"`
	Category string `json:"category" db:"category"`
}
