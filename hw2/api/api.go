package api

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo.Echo
}

func New() Server {
	return Server{
		*echo.New(),
	}
}

func (s *Server) Serve(port string) {
	s.HideBanner = true
	s.Logger.SetLevel(2)
	s.Use(middleware.Recover())
	s.GET("/", func(e echo.Context) error {
		s.Logger.Info("main logic is ok")
		return e.String(http.StatusOK, "Calm down and go on working")
	})
	s.GET("/itsalive", func(e echo.Context) error {
		s.Logger.Info("liveness is ok")
		return e.NoContent(http.StatusOK)
	})
	log.Fatal(s.Start(port))
}
