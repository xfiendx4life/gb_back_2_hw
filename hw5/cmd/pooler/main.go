package main

import (
	"fmt"
	"time"

	"github.com/xfiendx4life/gb_back_2_hw/hw5/activities"
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
	acts := []*activities.Activity{
		{UserId: 1, Name: "election", Date: time.Date(2020, time.November,
			3, 8, 0, 0, 0, time.FixedZone(time.UTC.String(), -8))},
		{UserId: 10, Name: "wait at home", Date: time.Date(2020, time.November,
			3, 8, 0, 0, 0, time.FixedZone(time.UTC.String(), -8))},
		{UserId: 13, Name: "2021 United States Capitol attack", Date: time.Date(2021, time.January,
			6, 10, 0, 0, 0, time.FixedZone(time.UTC.String(), -8))},
		{UserId: 25, Name: "have no idea", Date: time.Date(2021, time.January,
			6, 10, 0, 0, 0, time.FixedZone(time.UTC.String(), -8))},
	}

	for _, a := range acts {
		err := a.Create(m, p)
		if err != nil {
			fmt.Println(fmt.Errorf("error on create user %v: %w", a, err))
		}
	}
}
