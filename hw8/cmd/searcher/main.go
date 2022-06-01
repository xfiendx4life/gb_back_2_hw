package main

import (
	// "os"

	"os"

	"github.com/xfiendx4life/gb_back_2_hw/hw8/pkg/api"
)

func main() {
	addr := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")
	api.Start(addr, port)
}
