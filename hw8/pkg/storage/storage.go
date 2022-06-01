package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/google/uuid"
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
		data.Id = uuid.New()
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

func (s *Storage) Find(ctx context.Context, searchString string) ([](*models.Item), error) {
	query := fmt.Sprintf(`{
		"query": {
		  "query_string": {
			"query": "%s",
			"fields": [
			  "name^2",
			  "seller"
			],
			"type": "most_fields"
		  }
		}
	  }`, searchString)
	res, err := s.es.Search(
		s.es.Search.WithBody(strings.NewReader(query)),
		s.es.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("can't get data from es %s", err)
	}
	resData := make([](*models.Item), 0)
	var resMap map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&resMap)
	if err != nil {
		return nil, fmt.Errorf("can't decode result %s", err)
	}
	log.Println(resMap)
	for _, hit := range resMap["hits"].(map[string]interface{})["hits"].([]interface{}) {
		// log.Println(hit.(map[string]interface{})["_source"])
		m, err := models.MapToItem(hit.(map[string]interface{})["_source"].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		resData = append(resData, m)
	}
	return resData, nil
}
