package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
	rpc_server "github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/grpc/pricerdr/server"
	"github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/rest/prcr"
	"github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/rest/prcr/deliver"
	"github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/storage"
)

func main() {
	mode := flag.String("mode", "", "set up mode, use n to set a new dir to save data")
	serverToChoose := flag.String("server", "", "choose server to start. rest, rpc or empty for both")
	flag.Parse()
	dir := os.Getenv("STORAGE")
	st, err := storage.New(dir, *mode)
	if err != nil {
		log.Fatalf("can't set up storage %s", err)
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	switch {
	case *serverToChoose == "rpc":
		rpc := startRPC(ctx, st)
		<-sig
		rpc.Shutdown()
	case *serverToChoose == "rest":
		e := startREST(st)
		<-sig
		err = e.Shutdown(ctx)
		if err != nil {
			log.Fatalf("can't stop rest server: %s", err)
		}
	default:
		rpc := startRPC(ctx, st)
		e := startREST(st)
		<-sig
		rpc.Shutdown()
		err = e.Shutdown(ctx)
		if err != nil {
			log.Fatalf("can't stop rest server: %s", err)
		}
	}

}

func startRPC(ctx context.Context, st storage.Storage) *rpc_server.Server {
	server := rpc_server.New(st)
	var prt string
	if prt = os.Getenv("RPCPORT"); prt == "" {
		prt = ":8080"
	}
	go func() {
		rpc_server.Listen(ctx, server, prt)
	}()
	return server
}

func startREST(st storage.Storage) *echo.Echo {
	serv := echo.New()
	serv.HideBanner = true
	prcr.RegisterHandlers(serv, deliver.New(st))
	var prt string
	if prt = os.Getenv("RESTPORT"); prt == "" {
		prt = ":9000"
	}
	go serv.Start(prt)
	return serv
}
