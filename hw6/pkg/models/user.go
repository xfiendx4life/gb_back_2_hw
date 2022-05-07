package models

import "encoding/json"

type User struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	Confirmed bool   `json:"confirmed"`
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
