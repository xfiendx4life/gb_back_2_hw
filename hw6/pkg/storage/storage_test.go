package storage_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/models"
	"github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/storage"
)

var testClient = redis.NewClient(&redis.Options{
	Addr:     fmt.Sprintf("%s:%s", "localhost", "6379"),
	Password: "", // TODO: Add from config
	DB:       0,  // use default DB
})

func TestCreate(t *testing.T) {
	client, err := storage.NewRedisClient(
		"localhost",
		"6379",
		time.Duration(2*time.Minute),
	)
	assert.NoError(t, err)
	target := models.User{
		ID:       1,
		Name:     "testname",
		Password: "123",
	}
	err = client.Create(context.Background(), &target)
	assert.NoError(t, err)
	real, err := testClient.Get(context.Background(), "1").Bytes()
	assert.NoError(t, err)
	uReal := models.User{}
	_ = json.Unmarshal(real, &uReal)
	assert.Equal(t, target, uReal)
}
