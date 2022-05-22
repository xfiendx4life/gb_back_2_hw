package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/mediocregopher/radix/v3"
	"github.com/nats-io/nats.go"
	_ "gocloud.dev/pubsub/kafkapubsub"
)

var (
	connFunc = func(network, addr string) (radix.Conn, error) {
		return radix.Dial(network, addr,
			radix.DialTimeout(10*time.Second),
		)
	}
)

var brokerAddress = os.Getenv("BROKER_ADDRESS")
var topic = os.Getenv("TOPIC")

func main() {
	if brokerAddress == "" {
		brokerAddress = nats.DefaultURL
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	nc, err := nats.Connect(brokerAddress)
	if err != nil {
		log.Printf("can't connect to nats %s \n", err)
	}
	nc.Subscribe(topic, func(m *nats.Msg) {
		fmt.Printf("received: %s\n", string(m.Data))
		err = storage().Do(radix.Cmd(nil, "LPUSH", "result", string(m.Data)))
		if err != nil {
			log.Println(err)
		}
		if rand.Float64() < .05 {
			_ = storage().Do(radix.Cmd(nil, "LTRIM", "result", "0", "9"))
		}
	})
	<-ctx.Done()
	cancel()
}

func storage() *radix.Pool {
	var err error
	addr := os.Getenv("REDIS")
	s, err := radix.NewPool("tcp", addr, 1, radix.PoolConnFunc(connFunc))
	if err != nil {
		log.Panicf("can't connect to redis at %s addr %s", addr, err)
	}
	if s != nil {
		return s
	}
	return s
}
