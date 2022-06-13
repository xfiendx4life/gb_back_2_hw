package models

import (
	"github.com/google/uuid"
)

type Item struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

type List struct {
	ID    uuid.UUID `json:"id"`
	Items []*Item   `json:"items"`
}
