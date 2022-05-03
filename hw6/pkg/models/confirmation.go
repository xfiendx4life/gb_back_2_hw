package models

// TODO: Structure http://localhost:8080/confirmation/userId/hashedcode
type Confirmation struct {
	ID     int
	UserID int
	Code   int
}
