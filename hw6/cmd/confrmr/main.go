package main

import (
	"os"
	"strconv"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/api"
)

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	redisHost := os.Getenv("REDISHOST")
	redisPort := os.Getenv("REDISPORT")
	ttl, err := strconv.Atoi(os.Getenv("TTL"))
	if err != nil {
		log.Errorf("can't parse ttl %s", err)
		ttl = 1
	}
	server := api.New()
	server.Serve(host, port, redisHost, redisPort, time.Duration(ttl)*time.Hour)
}
