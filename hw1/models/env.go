package models

import "github.com/google/uuid"

type Env struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
