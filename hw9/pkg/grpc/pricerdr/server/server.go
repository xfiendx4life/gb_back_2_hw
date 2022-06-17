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

// TODO: Create Errors with status

type Server struct {
	*pricerdr.UnimplementedListServiceServer
	st   storage.Storage
	grpc *grpc.Server
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
	return &models.List{
		ID:    uuid.New(),
		Items: parseItemToModel(price.Items),
	}, nil
}
func parseItemsFromModel(list []*models.Item) []*pricerdr.Item {
	res := make([]*pricerdr.Item, len(list))
	for i, itm := range list {
		res[i] = &pricerdr.Item{}
		res[i].Name = itm.Name
		res[i].Price = itm.Price
	}
	return res
}
func parseListFromModel(list *models.List) *pricerdr.List {
	return &pricerdr.List{
		Id:    list.ID.String(),
		Items: parseItemsFromModel(list.Items),
	}
}

//* Creates pricelist
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

// * Reads data from price list and returns it
func (s *Server) Read(ctx context.Context, id *pricerdr.ListId) (*pricerdr.List, error) {
	select {
	case <-ctx.Done():
		return &pricerdr.List{}, fmt.Errorf("done with context")
	default:
		log.Printf("ready to read list with id %s", id.Id)
		uid, err := uuid.Parse(id.Id)
		if err != nil {
			return &pricerdr.List{}, fmt.Errorf("can't parse id to uuid %s", err)
		}
		list, err := s.st.Read(ctx, uid)
		if err != nil {
			return &pricerdr.List{}, fmt.Errorf("can't get data from storage: %s", err)
		}
		return parseListFromModel(list), nil
	}
}

// * An incoming priceList contents only rows to update
func (s *Server) Update(ctx context.Context, price *pricerdr.List) (*pricerdr.ListId, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("done with context")
	default:
		id, err := uuid.Parse(price.Id)
		if err != nil {
			return nil, fmt.Errorf("can't parse id: %s", err)
		}
		err = s.st.Update(ctx, id, parseItemToModel(price.Items))
		if err != nil {
			return nil, fmt.Errorf("can't update price %s: %s", price.Id, err)
		}
		return &pricerdr.ListId{Id: price.Id}, nil
	}
}

// * Deletes the whole pricelist
func (s *Server) Delete(ctx context.Context, price *pricerdr.ListId) (*pricerdr.ListId, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("done with context")
	default:
		uid, err := uuid.Parse(price.Id)
		if err != nil {
			return nil, fmt.Errorf("can't parse id to uuid %s", err)
		}
		err = s.st.Delete(ctx, uid)
		if err != nil {
			return nil, fmt.Errorf("can't delete pricelist %s", err)
		}
		return price, nil
	}
}

// * Listen registers server and starting it
func Listen(ctx context.Context, server *Server, prt string) error {
	server.grpc = grpc.NewServer()
	lis, err := net.Listen("tcp", prt)
	if err != nil {
		return fmt.Errorf("can't start listener")
	}
	log.Printf("listen to %s", prt)

	pricerdr.RegisterListServiceServer(server.grpc, server)
	reflection.Register(server.grpc)
	err = server.grpc.Serve(lis)
	if err != nil {
		return fmt.Errorf("can't serve: %s", err)
	}
	return nil
}

//* Gracefully shuts down server
func (s *Server) Shutdown() {
	s.grpc.GracefulStop()
	log.Println("server stopped gracefully")
}
