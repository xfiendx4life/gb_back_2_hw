package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/xfiendx4life/gb_back_2_hw/hw8/pkg/models"
)

type Storage struct {
	es *elasticsearch.Client
}

func New() (*Storage, error) {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, fmt.Errorf("can't create connection %s", err)
	}
	return &Storage{
		es: es,
	}, nil
}

func (s *Storage) Insert(ctx context.Context, data *models.Item) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("done with context")
	default:
		jsonString, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("can't marshal json %s", err)
		}
		request := esapi.IndexRequest{Index: "items", DocumentID: data.Id.String(), Body: strings.NewReader(string(jsonString))}
		resp, err := request.Do(ctx, s.es)
		_ = resp
		if err != nil {
			return fmt.Errorf("can't add data %s", err)
		}
	}
	return nil

}
