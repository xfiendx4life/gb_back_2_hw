package storage_test

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_back_2_hw/hw8/pkg/models"
	"github.com/xfiendx4life/gb_back_2_hw/hw8/pkg/storage"
)

func TestCheckConn(t *testing.T) {
	es, err := elasticsearch.NewDefaultClient()
	assert.NoError(t, err)

	res, err := es.Info()
	assert.NoError(t, err)

	defer res.Body.Close()
	log.Println(res)
}

func TestInsert(t *testing.T) {
	st, err := storage.New()
	assert.NoError(t, err)
	testData := &models.Item{
		Name:   "testname",
		Id:     uuid.New(),
		Price:  400,
		Seller: "testseller",
	}
	err = st.Insert(context.Background(), testData)
	assert.NoError(t, err)
	req := esapi.GetRequest{Index: "items", DocumentID: testData.Id.String()}
	es, err := elasticsearch.NewDefaultClient()
	assert.NoError(t, err)
	resp, err := req.Do(context.Background(), es)
	assert.NoError(t, err)
	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, testData.Name, res["_source"].(map[string]interface{})["name"].(string))
}
