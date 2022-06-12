package models

import (
	"github.com/google/uuid"
)

type Item struct {
	Name  string
	Price int64
}

type List struct {
	ID    uuid.UUID
	Items []*Item
}
