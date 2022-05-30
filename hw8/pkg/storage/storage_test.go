package storage_test

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_back_2_hw/hw8/pkg/models"
	"github.com/xfiendx4life/gb_back_2_hw/hw8/pkg/storage"
)

var cases = []models.Item{
	{
		Id:     uuid.New(),
		Name:   "test1",
		Price:  10,
		Seller: "testseller",
	},
	{
		Id:     uuid.New(),
		Name:   "test2",
		Price:  101,
		Seller: "test1",
	},
	{
		Id:     uuid.New(),
		Name:   "testseller",
		Price:  110,
		Seller: "testseller",
	},
}

var client = &http.Client{}

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
	defer resp.Body.Close()
	defer func() {
		req, err := http.NewRequest("DELETE", "http://localhost:9200/items", nil)
		assert.NoError(t, err)
		_, err = client.Do(req)
		assert.NoError(t, err)
	}()
	assert.NoError(t, err)
	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, testData.Name, res["_source"].(map[string]interface{})["name"].(string))
}

// TODO: Test Find
func TestFind(t *testing.T) {
	defer func() {
		req, err := http.NewRequest("DELETE", "http://localhost:9200/items", nil)
		assert.NoError(t, err)
		_, err = client.Do(req)
		assert.NoError(t, err)
	}()
	// ! I know i can't use func from previous test
	// ! But it's not interesting to copy same code here
	st, err := storage.New()
	assert.NoError(t, err)
	for _, item := range cases {
		err = st.Insert(context.Background(), &item)
		assert.NoError(t, err)
	}
	time.Sleep(time.Second)
	res, err := st.Find(context.Background(), "test1")

	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
}
