package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"

	"github.com/mediocregopher/radix/v3"
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
	topic          = "rates"
	broker1Address = "localhost:9092"
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
		HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Here we are"))
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
	wr := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker1Address},
		Topic:   topic,
	})
	if _, err := strconv.Atoi(rate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := wr.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(rate),
	})

	if err != nil {
		log.Printf("could not write message " + err.Error())
	}
}

// func topic() *pubsub.Topic {
// 	var err error
// 	t, err := pubsub.OpenTopic(context.Background(), "kafka://:9092/rates")
// 	if err != nil {
// 		panic(err)
// 	}
// 	if t != nil {
// 		return t
// 	}
// 	return t
// }

func storage() *radix.Pool {
	var err error
	s, err := radix.NewPool("tcp", ":6379", 1, radix.PoolConnFunc(connFunc))
	if err != nil {
		panic(err)
	}
	if s != nil {
		return s
	}
	return s
}
