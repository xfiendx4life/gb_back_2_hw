package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/xfiendx4life/gb_back_2_hw/hw8/pkg/models"
	"github.com/xfiendx4life/gb_back_2_hw/hw8/pkg/storage"
)

type Server struct {
	*http.Server
}

func New(port, addr string, router http.Handler) *Server {
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

func AddCard(w http.ResponseWriter, r *http.Request) {
	log.Println("Ready to add card")
	data := r.Body
	var card models.Item
	err := json.NewDecoder(data).Decode(&card)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	st, err := storage.New()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
	}
	w.WriteHeader(http.StatusAccepted)
}
