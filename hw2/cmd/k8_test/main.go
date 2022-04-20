package main

import (
	"log"
	"os"

	"github.com/xfiendx4life/gb_back_2_hw/hw2/api"
)

func main() {
	s := api.New()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("no port provided")
	}
	s.Serve(port)
}
