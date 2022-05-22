package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"

	"github.com/mediocregopher/radix/v3"
	"github.com/xfiendx4life/gb_back_2_hw/hw7/process"
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
	topic = "rates"
)

var (
	brokerAddress = os.Getenv("BROKER_ADDRESS")
	natsConn      *nats.Conn
	opened        int64
)

type Server struct {
	http.Server
}

func NewServer(addr string, router http.Handler) *Server {
	return &Server{
		http.Server{
			Addr:         addr,
			Handler:      router,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) Serve() {
	go s.ListenAndServe()
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	s.Shutdown(ctx)
	cancel()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	router := mux.NewRouter()
	router.
		HandleFunc("/rate", PostRateHandler).
		Methods(http.MethodPost)
	router.
		HandleFunc("/total", GetTotalHandler).
		Methods(http.MethodGet)
	router.
		HandleFunc("/opened", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(fmt.Sprintf("currently opened connections %d", opened)))
		})
	addr := os.Getenv("ADDRESS")
	s := NewServer(addr, router)
	s.Serve()
	log.Printf("serving at %s\n", addr)
	<-ctx.Done()
	s.Stop()
	cancel()
	log.Printf("stopped at %v", time.Now())
}

func GetTotalHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("access total handler")
	var rates []string
	err := storage().Do(radix.Cmd(&rates, "LRANGE", "result", "0", "10"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(rates) == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	var sum int
	for _, rate := range rates {
		v, err := strconv.Atoi(rate)
		if err != nil {
			continue
		}
		sum += v
	}
	result := float64(sum) / float64(len(rates))
	_, _ = w.Write([]byte(fmt.Sprintf("%.2f", result)))
}

func PostRateHandler(w http.ResponseWriter, r *http.Request) {

	rate := r.FormValue("rate")
	if _, err := strconv.Atoi(rate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	go func() {
		atomic.AddInt64(&opened, 1)
		process.New().Proceed(&opened)
	}()
	if brokerAddress == "" {
		brokerAddress = nats.DefaultURL
	}
	var err error
	if natsConn == nil {
		natsConn, err = nats.Connect(brokerAddress)
		if err != nil {
			log.Panicf("can't connect to nats on %s %s", brokerAddress, err)
		}

		if err != nil {
			log.Printf("could not write message " + err.Error())
		}
	}
	err = natsConn.Publish(topic, []byte(rate))
	if err != nil {
		log.Printf("can't publish to nats %s", err)
	}
	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	// <-ctx.Done()
	// cancel()

}

func storage() *radix.Pool {
	var err error
	addr := os.Getenv("REDIS")
	s, err := radix.NewPool("tcp", addr, 1, radix.PoolConnFunc(connFunc))
	if err != nil {
		panic(err)
	}
	if s != nil {
		return s
	}
	return s
}
