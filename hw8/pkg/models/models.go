package models

import "github.com/google/uuid"

type Item struct {
	Id     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Price  int64     `json:"price"`
	Seller string    `json:"seller"`
}
