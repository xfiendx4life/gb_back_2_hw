package models

type Student struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Faculty  string `json:"faculty"`
}
