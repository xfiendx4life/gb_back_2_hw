package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/mediocregopher/radix/v3"
	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/kafkapubsub"
)

var (
	connFunc = func(network, addr string) (radix.Conn, error) {
		return radix.Dial(network, addr,
			radix.DialTimeout(10*time.Second),
		)
	}
)

func subscription() (*pubsub.Subscription, error) {
	log.Println("enter in subscription")
	var err error
	sub, err := pubsub.OpenSubscription(context.Background(),
		"kafka://process?topic=rates")

	if err != nil {
		log.Println("error on open subscription")
		return nil, err
	}
	if sub != nil {
		return sub, nil
	}
	return sub, nil
}

func main() {
	for {
		s, err := subscription()
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			continue
		}
		msg, err := s.Receive(context.Background())
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			continue
		}
		err = storage().Do(radix.Cmd(nil, "LPUSH", "result", string(msg.Body)))
		if err != nil {
			log.Println(err)
		}
		if rand.Float64() < .05 {
			_ = storage().Do(radix.Cmd(nil, "LTRIM", "result", "0", "9"))
		}
		msg.Ack()
	}
}

func storage() *radix.Pool {
	var err error
	s, err := radix.NewPool("tcp", "redis:6379", 1, radix.PoolConnFunc(connFunc))
	if err != nil {
		panic(err)
	}
	if s != nil {
		return s
	}
	return s
}
