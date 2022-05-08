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
	"github.com/xfiendx4life/gb_back_2_hw/hw6/pkg/user/storage"
)

var (
	testClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", "localhost", "6379"),
		Password: "",
		DB:       0, // use default DB
	})
	target = models.User{
		Name:     "testname",
		Password: "123",
	}
	ttl = time.Duration(2 * time.Minute)
)

func TestCreate(t *testing.T) {
	client, err := storage.NewUserStorage(
		"localhost",
		"6379",
		time.Duration(2*time.Minute),
	)
	assert.NoError(t, err)
	err = client.Create(context.Background(), &target)
	assert.NoError(t, err)
	real, err := testClient.Get(context.Background(), "testname").Bytes()
	assert.NoError(t, err)
	uReal := models.User{}
	_ = json.Unmarshal(real, &uReal)
	assert.Equal(t, target, uReal)
	testClient.Del(context.Background(), "testname")
}

func TestCreateOnExistingKey(t *testing.T) {
	client, err := storage.NewUserStorage(
		"localhost",
		"6379",
		time.Duration(2*time.Minute),
	)
	assert.NoError(t, err)
	err = client.Create(context.Background(), &target)
	assert.NoError(t, err)
	tt := target
	tt.Confirmed = true
	err = client.Create(context.Background(), &tt)
	assert.NoError(t, err)
	real, err := testClient.Get(context.Background(), "testname").Bytes()
	assert.NoError(t, err)
	uReal := models.User{}
	_ = json.Unmarshal(real, &uReal)
	assert.True(t, uReal.Confirmed)
	testClient.Del(context.Background(), "testname")
}

func TestGet(t *testing.T) {
	testClient.Set(context.Background(), "testname", &target, ttl)
	client, err := storage.NewUserStorage(
		"localhost",
		"6379",
		time.Duration(2*time.Minute),
	)
	assert.NoError(t, err)
	res, err := client.GetUser(context.Background(), "testname")
	assert.NoError(t, err)
	assert.Equal(t, target, *res)
}
