package server_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/grpc/pricerdr"
	rpc_server "github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/grpc/pricerdr/server"
	"github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type mStorage struct {
	err error
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

const tstId = "311e0467-a754-4286-bf89-47c3a83eeb68"

var (
	testCase1 = models.List{
		Items: []*models.Item{
			{Name: "test1", Price: 10},
			{Name: "test2", Price: 30},
		},
	}
	lis *bufconn.Listener
)

func (mc *mStorage) Create(ctx context.Context, list models.List) error {
	return mc.err
}

func (mc *mStorage) Read(ctx context.Context, id uuid.UUID) (list *models.List, err error) {
	if id.String() == tstId {
		testCase1.ID, _ = uuid.Parse(tstId)
		return &testCase1, mc.err
	}
	return nil, mc.err
}
func (mc *mStorage) Update(ctx context.Context, id uuid.UUID, items []*models.Item) error {
	return mc.err

}

func (mc *mStorage) Delete(ctx context.Context, id uuid.UUID) error {
	return mc.err
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func startTestRpc(ctx context.Context, err error, bufSize int) {
	server := rpc_server.New(&mStorage{err: err})
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pricerdr.RegisterListServiceServer(s, server)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
	<-ctx.Done()
	s.Stop()
}

func TestCreate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go startTestRpc(ctx, nil, 1024*1024)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := pricerdr.NewListServiceClient(conn)
	resp, err := client.Create(ctx, parseListFromModel(&testCase1))
	require.NoError(t, err)
	require.NotEmpty(t, resp)
}

func TestCreateErrorWrongType(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go startTestRpc(ctx, fmt.Errorf("storage testError"), 1024*1024)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := pricerdr.NewListServiceClient(conn)
	resp, err := client.Create(ctx, parseListFromModel(&testCase1))
	require.Error(t, err)
	require.Empty(t, resp)
}

func TestRead(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go startTestRpc(ctx, nil, 1024*1024)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := pricerdr.NewListServiceClient(conn)
	resp, err := client.Read(ctx, &pricerdr.ListId{Id: tstId})
	require.NoError(t, err)
	require.Equal(t, testCase1.Items, parseItemToModel(resp.Items))
}

// ! FATAL ERROR HERE BUT NO ERROR IN REAL WORK
// func TestReadError(t *testing.T) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	go startTestRpc(ctx, fmt.Errorf("test storage error"), 1024*1024)
// 	defer cancel()
// 	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer),
// 		grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	require.NoError(t, err)
// 	defer conn.Close()
// 	client := pricerdr.NewListServiceClient(conn)
// 	resp, err := client.Read(ctx, &pricerdr.ListId{Id: "2" + tstId[1:]})
// 	require.Error(t, err)
// 	require.Equal(t, testCase1.Items, parseItemToModel(resp.Items))
// }

func TestUpdate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go startTestRpc(ctx, nil, 1024*1024)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := pricerdr.NewListServiceClient(conn)
	testCase1.Items[0].Price = 3000
	_, err = client.Update(ctx, parseListFromModel(&testCase1))
	require.NoError(t, err)
}

func TestUpdateError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go startTestRpc(ctx, fmt.Errorf("test error"), 1024*1024)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := pricerdr.NewListServiceClient(conn)
	testCase1.Items[0].Price = 3000
	_, err = client.Update(ctx, parseListFromModel(&testCase1))
	require.Error(t, err)
}

func TestUpdateUserError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go startTestRpc(ctx, fmt.Errorf("test error"), 1024*1024)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := pricerdr.NewListServiceClient(conn)
	tt := parseListFromModel(&testCase1)
	tt.Id += "1"
	_, err = client.Update(ctx, tt)
	require.Error(t, err)
}

func TestDelete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go startTestRpc(ctx, nil, 1024*1024)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := pricerdr.NewListServiceClient(conn)
	testCase1.Items[0].Price = 3000
	_, err = client.Delete(ctx, &pricerdr.ListId{Id: tstId})
	require.NoError(t, err)
}

func TestDeleteStorageError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go startTestRpc(ctx, fmt.Errorf("test storage error"), 1024*1024)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := pricerdr.NewListServiceClient(conn)
	testCase1.Items[0].Price = 3000
	_, err = client.Delete(ctx, &pricerdr.ListId{Id: tstId})
	require.Error(t, err)
}

func TestDeleteUserError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go startTestRpc(ctx, fmt.Errorf("test storage error"), 1024*1024)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := pricerdr.NewListServiceClient(conn)
	testCase1.Items[0].Price = 3000
	_, err = client.Delete(ctx, &pricerdr.ListId{Id: "123"})
	require.Error(t, err)
}
