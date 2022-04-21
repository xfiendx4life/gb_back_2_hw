package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xfiendx4life/gb_back_2_hw/hw4/deliver"
	"github.com/xfiendx4life/gb_back_2_hw/hw4/metrics"
	"github.com/xfiendx4life/gb_back_2_hw/hw4/storage"
)

func main() {
	http.Handle("/metrics", promhttp.Handler())
	server := echo.New()
	server.HideBanner = true
	st, err := storage.New("storage.sqlite", metrics.New(true))
	if err != nil {
		log.Fatal(err)
	}
	del := deliver.New(st)
	server.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	server.POST("/add", del.SetStudent)
	server.GET("/get", del.GetByLastName)
	server.GET("/byfaculty", del.GetAllByFaculty)
	go func() {
		err := server.Start(":8080")
		if err != nil {
			log.Fatalf("can't start server %s", err)
		}
	}()
	_, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("server stopping")
	cancel()
}
