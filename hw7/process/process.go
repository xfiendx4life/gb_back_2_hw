package process

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync/atomic"
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

const (
	topic = "rates" // ? make it a New() parameter
)

type Process struct {
	brokerAddress string
}

func New() *Process {
	brokerAddress := os.Getenv("BROKER_ADDRESS")
	if brokerAddress == "" {
		brokerAddress = nats.DefaultURL
	}
	return &Process{
		brokerAddress: brokerAddress,
	}
}

func (p *Process) Proceed(counter *int64) {
	nc, err := nats.Connect(p.brokerAddress)
	if err != nil {
		log.Printf("can't connect to nats %s \n", err)
	}
	nc.Subscribe(topic, func(m *nats.Msg) {
		fmt.Printf("received: %s\n", string(m.Data))
		pool := storage()
		err = pool.Do(radix.Cmd(nil, "LPUSH", "result", string(m.Data)))
		if err != nil {
			log.Println(err)
		}

		if rand.Float64() < .05 {
			_ = pool.Do(radix.Cmd(nil, "LTRIM", "result", "0", "9"))
		}
		pool.Close()
		nc.Close()
		atomic.AddInt64(counter, -1)
	})

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
