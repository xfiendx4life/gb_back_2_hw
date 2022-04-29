package main

import (
	"fmt"
	"log"
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
		Role:    "master",
		Address: "port=8100 user=test password=test dbname=test sslmode=disable",
		Number:  0,
	})
	m.Add(&manager.Shard{
		Role:    "master",
		Address: "port=8110 user=test password=test dbname=test sslmode=disable",
		Number:  1,
	})
	m.Add(&manager.Shard{
		Role:    "master",
		Address: "port=8120 user=test password=test dbname=test sslmode=disable",
		Number:  2,
	})
	m.Add(&manager.Shard{
		Role:    "slave",
		Address: "port=8101 user=test password=test dbname=test sslmode=disable",
		Number:  10,
	})
	m.Add(&manager.Shard{
		Role:    "slave",
		Address: "port=8111 user=test password=test dbname=test sslmode=disable",
		Number:  11,
	})
	m.Add(&manager.Shard{
		Role:    "slave",
		Address: "port=8121 user=test password=test dbname=test sslmode=disable",
		Number:  12,
	})
	fmt.Println("-----Creating users-----")
	uu := []*user.User{
		{UserId: 1, Name: "Joe Biden", Age: 78, Spouse: 10},
		{UserId: 11, Name: "Jill Biden", Age: 69, Spouse: 1},
		{UserId: 15, Name: "Donald Trump", Age: 74, Spouse: 25},
		{UserId: 26, Name: "Melania Trump", Age: 52, Spouse: 13},
	}
	for _, u := range uu {
		err := u.Create(m, p)
		if err != nil {
			fmt.Println(fmt.Errorf("error on create user %v: %w", u, err))
		}
	}
	fmt.Println("-----Creating activities-----")
	acts := []*activities.Activity{
		{UserId: 1, Name: "election", Date: time.Date(2020, time.November,
			3, 8, 0, 0, 0, time.FixedZone(time.UTC.String(), -8))},
		{UserId: 11, Name: "wait at home", Date: time.Date(2020, time.November,
			3, 8, 0, 0, 0, time.FixedZone(time.UTC.String(), -8))},
		{UserId: 15, Name: "2021 United States Capitol attack", Date: time.Date(2021, time.January,
			6, 10, 0, 0, 0, time.FixedZone(time.UTC.String(), -8))},
		{UserId: 26, Name: "have no idea", Date: time.Date(2021, time.January,
			6, 10, 0, 0, 0, time.FixedZone(time.UTC.String(), -8))},
	}

	for _, a := range acts {
		err := a.Create(m, p)
		if err != nil {
			fmt.Println(fmt.Errorf("error on create activity %v: %w", a, err))
		}
	}

	fmt.Println("-----Reading users-----")
	usersToRead := []user.User{
		{UserId: 11},
		{UserId: 1},
		{UserId: 15},
		{UserId: 26},
	}
	for _, u := range usersToRead {
		err := u.Read(m, p)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(u)
		}
	}
	fmt.Println("-----Reading activities-----")
	actsToRead := []activities.Activity{
		{UserId: 11},
		{UserId: 1},
		{UserId: 15},
		{UserId: 26},
	}
	for _, a := range actsToRead {
		err := a.Read(m, p)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(a)
		}
	}

	fmt.Println("-----Deleting users-----")

	for _, u := range uu {
		err := u.Delete(m, p)
		if err != nil {
			fmt.Println(fmt.Errorf("error on delete user %v: %w", u, err))
		}
	}
	fmt.Println("-----Deleting activities-----")
	for _, a := range acts {
		err := a.Delete(m, p)
		if err != nil {
			fmt.Println(fmt.Errorf("error on delete activity %v: %w", a, err))
		}
	}

}
