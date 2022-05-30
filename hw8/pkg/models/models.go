package models

import (
	"github.com/google/uuid"
)

type Item struct {
	Id     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Price  int64     `json:"price"`
	Seller string    `json:"seller"`
}

func MapToItem(data map[string]interface{}) (m *Item, err error) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		log.Printf("%v", r)
	// 		log.Println("panic recovered")
	// 		err = fmt.Errorf("panic occured, can't parse data to Item model")
	// 	}
	// }()
	m = &Item{}
	m.Id, err = uuid.Parse(data["id"].(string))
	m.Name = data["name"].(string)
	m.Price = int64(data["price"].(float64))
	m.Seller = data["seller"].(string)

	return
}
