package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/delivery"
)

type Server struct {
	echo.Echo
}

func New() *Server {
	return &Server{
		Echo: *echo.New(),
	}
}

func (s *Server) Serve(host, port string, redisHost, redisPort string, ttl time.Duration) {
	d, err := delivery.New(redisHost, redisPort, ttl)
	if err != nil {
		s.Logger.Errorf("can't initialized server %s", err)
	}
	s.HideBanner = true
	s.Use(middleware.Recover())
	s.POST("/", d.CreateUser)
	s.POST("/:user/confirm", d.Confirm)
	go func() {
		h := fmt.Sprintf("%s:%s", host, port)
		if err := s.Start(h); err != nil && err != http.ErrServerClosed {
			s.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		s.Logger.Fatal(err)
	}
}
