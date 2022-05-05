package models

import "encoding/json"

// TODO: Structure http://localhost:8080/confirmation/userId/hashedcode
type Confirmation struct {
	UserID int    `json:"userId"`
	Code   string `json:"code"`
}

func (c *Confirmation) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Confirmation) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}
