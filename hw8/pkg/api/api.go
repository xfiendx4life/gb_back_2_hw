package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/xfiendx4life/gb_back_2_hw/hw8/pkg/models"
	"github.com/xfiendx4life/gb_back_2_hw/hw8/pkg/storage"
)

type Server struct {
	*http.Server
}

func New(addr, port string, router http.Handler) *Server {
	return &Server{
		&http.Server{
			Addr:         fmt.Sprintf("%s:%s", addr, port),
			ReadTimeout:  time.Second * 2,
			WriteTimeout: time.Second * 2,
			Handler:      router,
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

func AddCard(w http.ResponseWriter, r *http.Request) {
	log.Println("Ready to add card")
	var card models.Item
	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		log.Printf("can't decode body %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	st, err := storage.New()
	if err != nil {
		log.Printf("can't connect to storage %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	doneCh := make(chan struct{})
	go func() {
		err = st.Insert(context.Background(), &card)
		doneCh <- struct{}{}
	}()
	<-doneCh
	if err != nil {
		log.Printf("error while adding data to strage %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("data accepted\n")
	w.WriteHeader(http.StatusAccepted)
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	log.Println("ready to get item")
	searchString, ok := mux.Vars(r)["search"]
	if !ok {
		log.Printf("Can't get vars from path \n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	st, err := storage.New()
	if err != nil {
		log.Printf("Error while connecting to storage %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	resChan := make(chan []*models.Item, 1)
	go func(ctx context.Context) {
		var res []*models.Item
		res, err = st.Find(ctx, searchString)
		resChan <- res
	}(ctx)
	select {
	case <-ctx.Done():
		log.Println("Done with context")
	case res := <-resChan:
		log.Println("Got answer")
		json.NewEncoder(w).Encode(res)
	}
	cancel()
}

func Start(addr, port string) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	router := mux.NewRouter()
	router.
		HandleFunc("/add", AddCard).
		Methods(http.MethodPost)
	router.
		HandleFunc("/item/{search}", GetItem).
		Methods(http.MethodGet)
	// addr := os.Getenv("ADDRESS")
	s := New(addr, port, router)
	s.Serve()
	log.Printf("serving at %s:%s\n", addr, port)
	<-ctx.Done()
	s.Stop()
	cancel()
	log.Printf("stopped at %v", time.Now())
}
