package api

import (
	"github.com/labstack/echo/v4"
	conDeliver "github.com/xfiendx4life/gb_back_2_hw/internal/pkg/connector/deliver"
	conCase "github.com/xfiendx4life/gb_back_2_hw/internal/pkg/connector/usecase"
	envDeliver "github.com/xfiendx4life/gb_back_2_hw/internal/pkg/env/deliver"
	envStorage "github.com/xfiendx4life/gb_back_2_hw/internal/pkg/env/storage"
	envCase "github.com/xfiendx4life/gb_back_2_hw/internal/pkg/env/usecase"
	userDeliver "github.com/xfiendx4life/gb_back_2_hw/internal/pkg/user/deliver"
	userStorage "github.com/xfiendx4life/gb_back_2_hw/internal/pkg/user/storage"
	userCase "github.com/xfiendx4life/gb_back_2_hw/internal/pkg/user/usecase"
)

type Server struct {
	s *echo.Echo
}

func New() Server {
	server := echo.New()
	server.HideBanner = true
	return Server{
		s: server,
	}
}

func (s *Server) Run() {
	uCase := userCase.New(userStorage.New())
	userDel := userDeliver.New(uCase)

	eCase := envCase.New(envStorage.New())
	envDel := envDeliver.New(eCase)

	connDel := conDeliver.New(conCase.New(uCase, eCase))

	s.s.GET("/user:name", userDel.GetByName)
	s.s.POST("/user/create", userDel.Create)
	s.s.GET("/env:name", envDel.Get)
	s.s.POST("/env/create", envDel.Create)
	s.s.POST("/user/add_to_env:user:env", connDel.AddToEnv)
	s.s.GET("/env:user", connDel.GetByUser)
	s.s.GET("/user:user", connDel.GetByEnv)
	s.s.DELETE("/user:env", connDel.DeleteUserFromEnv)

	s.s.Logger.Fatal(s.s.Start(":8000"))
}
