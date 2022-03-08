package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID     uuid.UUID
	Name   string
	Envs []Env  // TODO something more realistic
}
