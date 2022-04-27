package main

import (
	"fmt"

	"github.com/xfiendx4life/gb_back_2_hw/hw5/manager"
	"github.com/xfiendx4life/gb_back_2_hw/hw5/pool"
	"github.com/xfiendx4life/gb_back_2_hw/hw5/user"
)

func main() {
	m := manager.NewManager(3)
	p := pool.NewPool()
	m.Add(&manager.Shard{
		Address: "port=8100 user=test password=test dbname=test sslmode=disable",
		Number:  0,
	})
	m.Add(&manager.Shard{
		Address: "port=8110 user=test password=test dbname=test sslmode=disable",
		Number:  1,
	})
	m.Add(&manager.Shard{
		Address: "port=8120 user=test password=test dbname=test sslmode=disable",
		Number:  2,
	})
	uu := []*user.User{
		{UserId: 1, Name: "Joe Biden", Age: 78, Spouse: 10},
		{UserId: 10, Name: "Jill Biden", Age: 69, Spouse: 1},
		{UserId: 13, Name: "Donald Trump", Age: 74, Spouse: 25},
		{UserId: 25, Name: "Melania Trump", Age: 52, Spouse: 13},
	}
	for _, u := range uu {
		err := u.Create(m, p)
		if err != nil {
			fmt.Println(fmt.Errorf("error on create user %v: %w", u, err))
		}
	}
}
