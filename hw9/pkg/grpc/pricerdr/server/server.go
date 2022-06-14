package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/grpc/pricerdr"
	"github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/models"
	"github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	*pricerdr.UnimplementedListServiceServer
	st storage.Storage
}

func New(st storage.Storage) *Server {
	return &Server{
		st: st,
	}
}

func parseItemToModel(items []*pricerdr.Item) []*models.Item {
	log.Println(len(items))
	itms := make([]*models.Item, len(items))
	for i, itm := range items {
		log.Println(itm.Name, itm.Price)
		itms[i] = &models.Item{}
		itms[i].Name = itm.Name
		itms[i].Price = itm.Price
	}
	return itms
}

func parseListToModel(price *pricerdr.List) (*models.List, error) {
	//! ignore uuid, create new one
	// id, err := uuid.FromBytes([]byte(price.Id))
	// if err != nil {
	// 	log.Printf("can't parse uuid: %s", err)
	// 	return nil, err
	// }
	return &models.List{
		ID:    uuid.New(),
		Items: parseItemToModel(price.Items),
	}, nil
}

func (s *Server) Create(ctx context.Context, price *pricerdr.List) (*pricerdr.ListId, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("done with context")
	default:
		log.Printf("recieve pricelist for creating %v", price)
		data, err := parseListToModel(price)
		if err != nil {
			return nil, fmt.Errorf("can't parse input to models: %s", err)
		}
		err = s.st.Create(ctx, *data)
		if err != nil {
			return nil, fmt.Errorf("can't write data to storage: %s", err)
		}
		return &pricerdr.ListId{
			Id: data.ID.String(),
		}, nil

	}

}

func (s *Server) Read(context.Context, *pricerdr.ListId) (*pricerdr.List, error) {
	return nil, nil
}
func (s *Server) Update(context.Context, *pricerdr.List) (*pricerdr.ListId, error) {
	return nil, nil
}

func (s *Server) Delete(context.Context, *pricerdr.ListId) (*pricerdr.ListId, error) {
	return nil, nil
}

func Listen(server Server, prt string) error {
	lis, err := net.Listen("tcp", prt)
	if err != nil {
		return fmt.Errorf("can't start listener")
	}
	log.Printf("listen to %s", prt)
	grpcServer := grpc.NewServer()
	pricerdr.RegisterListServiceServer(grpcServer, &server)
	reflection.Register(grpcServer)
	err = grpcServer.Serve(lis)
	if err != nil {
		return fmt.Errorf("can't serve: %s", err)
	}
	return nil
}
