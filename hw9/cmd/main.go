package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	rpc_server "github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/grpc/pricerdr/server"
	"github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/storage"
)

func main() {
	mode := flag.String("mode", "", "set up mode, use n to set a new dir to save data")
	flag.Parse()
	dir := os.Getenv("STORAGE")
	st, err := storage.New(dir, *mode)
	if err != nil {
		log.Fatalf("can't set up storage %s", err)
	}
	server := rpc_server.New(st)
	var prt string
	if prt = os.Getenv("PORT"); prt == "" {
		prt = ":8080"
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	go func() {
		rpc_server.Listen(ctx, server, prt)
	}()
	<-ctx.Done()
	server.Shutdown()
}
