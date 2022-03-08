package models

import "github.com/google/uuid"

type Env struct {
	ID    uuid.UUID
	Name  string
	Users []User
}
