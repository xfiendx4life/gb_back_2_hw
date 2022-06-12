package server

import (
	"context"

	"github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/grpc/pricerdr"
)

type Server struct {
	*pricerdr.UnimplementedListServiceServer
}

func (s *Server) Create(context.Context, *pricerdr.List) (*pricerdr.ListId, error) {

}

func (s *Server) Read(context.Context, *pricerdr.ListId) (*pricerdr.List, error) {

}
func (s *Server) Update(context.Context, *pricerdr.List) (*pricerdr.ListId, error) {

}

func (s *Server) Delete(context.Context, *pricerdr.ListId) (*pricerdr.ListId, error) {

}
